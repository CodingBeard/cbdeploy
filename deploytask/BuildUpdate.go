package deploytask

import (
	"github.com/codingbeard/cbdeploy"
	"github.com/codingbeard/cbutil"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type FileDownload struct {
	RemotePath string
	LocalPath  string
}

type BuildUpdate struct {
	GetDownloader         func() (cbdeploy.Downloader, error)
	Log                   cbdeploy.Logger
	ErrorHandler          cbdeploy.ErrorHandler
	CurrentVersion        string
	VersionFileRemotePath string
	Bucket                string
	Files                 []FileDownload
	InitScriptRemotePath  string
}

func (c BuildUpdate) GetSchedule() string {
	return "manual"
}

func (c BuildUpdate) GetGroup() string {
	return "build"
}

func (c BuildUpdate) GetName() string {
	return "update"
}

func (c BuildUpdate) Run() error {
	var downloader cbdeploy.Downloader
	if c.GetDownloader != nil {
		var e error
		downloader, e = c.GetDownloader()
		if e != nil {
			c.ErrorHandler.Error(e)
			return e
		}
	}

	localInitScriptPath := "/usr/local/bin/init.sh"
	if c.InitScriptRemotePath != "" {
		c.Files = append(c.Files, FileDownload{c.InitScriptRemotePath, localInitScriptPath})
	}

	cbutil.RepeatingTask{
		Sleep:      time.Second * 5,
		SleepFirst: true,
		Blocking:   true,
		Run: func() {
			versionBytes, e := downloader.Download(c.Bucket, c.VersionFileRemotePath)
			if e != nil {
				c.ErrorHandler.Error(e)
				return
			}
			version := string(versionBytes)

			if version != c.CurrentVersion {
				c.Log.InfoF("root", "Different version found: %s, current version: %s", version, c.CurrentVersion)

				for _, file := range c.Files {
					c.Log.InfoF("DEPLOY", "Downloading file %s/%s", c.Bucket, file.RemotePath)
					initBytes, e := downloader.Download(c.Bucket, file.RemotePath)
					if e != nil {
						c.ErrorHandler.Error(e)
						return
					}

					c.Log.InfoF("DEPLOY", "Writing data to %s", file.LocalPath)
					e = ioutil.WriteFile(file.LocalPath, initBytes, os.ModePerm)
					if e != nil {
						c.ErrorHandler.Error(e)
						return
					}
				}

				if c.InitScriptRemotePath != "" {
					c.Log.InfoF("DEPLOY", "Making init script executable: %s", localInitScriptPath)
					e = exec.Command("chmod", "+x", localInitScriptPath).Run()
					if e != nil {
						c.ErrorHandler.Error(e)
						return
					}

					c.Log.InfoF("DEPLOY", "Executing init script: %s", localInitScriptPath)
					output, e := exec.Command(localInitScriptPath).Output()
					if e != nil {
						c.ErrorHandler.Error(e)
						return
					}
					c.Log.InfoF("root", string(output))
				}

				c.CurrentVersion = version
			}
		},
	}.Start()

	return nil
}
