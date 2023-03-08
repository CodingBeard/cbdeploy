package deploytask

import (
	"fmt"
	"github.com/codingbeard/cbdeploy"
)

type FileRelease struct {
	VersionPath string
	LivePath    string
}
type BuildRelease struct {
	ShouldRun     func() error
	GetDownloader func() (cbdeploy.Downloader, error)
	GetUploader   func() (cbdeploy.Uploader, error)
	Log           cbdeploy.Logger
	ErrorHandler  cbdeploy.ErrorHandler
	PublicUpload  bool
	Bucket        string
	Files         []FileRelease
}

func (c BuildRelease) GetSchedule() string {
	return "manual"
}

func (c BuildRelease) GetGroup() string {
	return "build"
}

func (c BuildRelease) GetName() string {
	return "release"
}

func (c BuildRelease) Run() error {
	if c.ShouldRun != nil {
		e := c.ShouldRun()
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
	}
	var uploader cbdeploy.Uploader
	var downloader cbdeploy.Downloader
	if c.GetUploader == nil {
		e := fmt.Errorf("BuildRelease.GetUploader == nil")
		c.ErrorHandler.Error(e)
		return e
	}

	if c.GetDownloader == nil {
		e := fmt.Errorf("BuildRelease.GetDownloader == nil")
		c.ErrorHandler.Error(e)
		return e
	}

	var e error
	uploader, e = c.GetUploader()
	if e != nil {
		c.ErrorHandler.Error(e)
		return e
	}
	if uploader == nil {
		e := fmt.Errorf("BuildRelease.GetUploader() == nil")
		c.ErrorHandler.Error(e)
		return e
	}

	downloader, e = c.GetDownloader()
	if e != nil {
		c.ErrorHandler.Error(e)
		return e
	}
	if downloader == nil {
		e := fmt.Errorf("BuildRelease.GetDownloader() == nil")
		c.ErrorHandler.Error(e)
		return e
	}

	for _, paths := range c.Files {
		c.Log.InfoF("DEPLOY", "Downloading file %s", paths.VersionPath)
		fileBytes, e := downloader.Download(c.Bucket, paths.VersionPath)
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
		c.Log.InfoF("DEPLOY", "Uploading file %s to %s/%s", paths.VersionPath, c.Bucket, paths.LivePath)
		e = uploader.UploadBytes(c.Bucket, paths.LivePath, fileBytes, c.PublicUpload)
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
	}

	return nil
}
