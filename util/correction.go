package util

import (
	"bytes"
	"github.com/krakowski/ilias"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"text/template"
)

const (
	CorrectionTemplateFilename = "CORRECTION.tmpl"
)

type CorrectionTemplate struct {
	Checksum	[20]byte
	Content		[]byte
}

type TemplateParams struct {
	Student		string
	Tutor 		string
}

func WriteCorrectionTemplate(path string, params TemplateParams) error {
	tpl, err := template.ParseFiles(CorrectionTemplateFilename)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	err = tpl.Execute(&out, params)
	if err != nil {
		return err
	}

	data := out.Bytes()
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func ReadCorrection(path string) (*ilias.Correction, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	correction := ilias.Correction{}
	err = yaml.Unmarshal(file, &correction)
	if err != nil {
		return nil, err
	}

	return &correction, nil
}

func FilterCorrections(values []ilias.Correction, test func(correction ilias.Correction) bool) (ret []ilias.Correction) {
	for _, s := range values {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
