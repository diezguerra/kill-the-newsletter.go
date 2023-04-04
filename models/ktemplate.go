package models

import (
	"bytes"
	htemplate "html/template"
	"io"
	"ktn-go/config"
	"strings"
	ttemplate "text/template"

	log "github.com/sirupsen/logrus"
)

type Executable interface {
	Execute(io.Writer, any) error
}

func RenderTemplate(data map[string]interface{}, files []string) ([]byte, error) {

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

	templData := make(map[string]interface{})
	for key, value := range data {
		templData[key] = value
	}
	for key, value := range config.TemplateVariables {
		templData[key] = value
	}

	executedTmpl := new(bytes.Buffer)
	if err := tmpl.(Executable).Execute(executedTmpl, templData); err != nil {
		log.Error("KTemplate: Couldn't execute template ", files, ": ", err)
		return nil, err
	}

	return executedTmpl.Bytes(), nil
}
