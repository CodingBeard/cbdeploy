package deploytask

import (
	"github.com/codingbeard/cbdeploy"
	"io/ioutil"
)

type FileUpload struct {
	LocalPath  string
	RemotePath string
}

type BuildDeploy struct {
	ShouldRun    func() error
	GetUploader  func() (cbdeploy.Uploader, error)
	Log          cbdeploy.Logger
	ErrorHandler cbdeploy.ErrorHandler
	PublicUpload bool
	Bucket       string
	Files        []FileUpload
}

func (c BuildDeploy) GetSchedule() string {
	return "manual"
}

func (c BuildDeploy) GetGroup() string {
	return "build"
}

func (c BuildDeploy) GetName() string {
	return "deploy"
}

func (c BuildDeploy) Run() error {
	if c.ShouldRun != nil {
		e := c.ShouldRun()
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
	}
	var uploader cbdeploy.Uploader
	if c.GetUploader != nil {
		var e error
		uploader, e = c.GetUploader()
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
	}

	for _, paths := range c.Files {
		c.Log.InfoF("DEPLOY", "Reading file %s", paths.LocalPath)
		fileBytes, e := ioutil.ReadFile(paths.LocalPath)
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
		c.Log.InfoF("DEPLOY", "Uploading file %s to %s/%s", paths.LocalPath, c.Bucket, paths.RemotePath)
		_ = fileBytes
		e = uploader.UploadBytes(c.Bucket, paths.RemotePath, fileBytes, c.PublicUpload)
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
	}

	return nil
}
