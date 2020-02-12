package cbdeploy

import (
	"fmt"
	"log"
)

type Logger interface {
	InfoF(category string, message string, args ...interface{})
	ErrorF(category string, message string, args ...interface{})
}

type DefaultLogger struct{}

func (d DefaultLogger) InfoF(category string, message string, args ...interface{}) {
	log.Println(category+":", fmt.Sprintf(message, args...))
}

func (d DefaultLogger) ErrorF(category string, message string, args ...interface{}) {
	log.Println(category+":", fmt.Sprintf(message, args...))
}

type ErrorHandler interface {
	Error(e error)
}

type DefaultErrorHandler struct{}

type Uploader interface {
	UploadBytes(bucket string, name string, data []byte, public bool) error
}

type Downloader interface {
	Download(bucket string, name string) ([]byte, error)
}
