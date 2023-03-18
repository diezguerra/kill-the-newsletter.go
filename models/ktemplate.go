package models

import (
	"bytes"
	htemplate "html/template"
	"io"
	"strings"
	ttemplate "text/template"

	log "github.com/sirupsen/logrus"
)

type Executable interface {
	Execute(io.Writer, any) error
}

func RenderTemplate(data any, files []string) ([]byte, error) {

	var tmpl interface{}
	var err error
	switch strings.HasSuffix(files[0], ".xml") {
	case true:
		tmpl, err = ttemplate.ParseFiles(files...)
	case false:
		tmpl, err = htemplate.ParseFiles(files...)
	}

	if err != nil {
		log.Error("KTemplate: couldn't load template ", files, ": ", err)
		return nil, err
	}

	executedTmpl := new(bytes.Buffer)
	if err := tmpl.(Executable).Execute(executedTmpl, data); err != nil {
		log.Error("KTemplate: Couldn't execute template ", files, ": ", err)
		return nil, err
	}

	return executedTmpl.Bytes(), nil
}
