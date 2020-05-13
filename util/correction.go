package util

import (
	"bytes"
	"github.com/krakowski/ilias"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const (
	CorrectionTemplateFilename = "CORRECTION.tmpl"
	SubmissionFilename = "Abgabe"
	CorrectionFilename = "Korrektur.yml"
)

type CorrectionTemplate struct {
	Checksum	[20]byte
	Content		[]byte
}

type TemplateParams struct {
	Student		string
	Tutor 		string
}

type CorrectionStats struct {
	Corrected	[]ilias.Correction
	Pending		[]ilias.Correction
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

func ReadCorrections(members []string) ([]ilias.Correction, error) {
	var corrections []ilias.Correction
	for _, member := range members {
		path := filepath.Join(member, CorrectionFilename)
		correction, err := ReadCorrection(path)
		if err != nil {
			return nil, err
		}

		corrections = append(corrections, *correction)
	}

	return corrections, nil
}

func GetCorrectionStats(corrections []ilias.Correction) CorrectionStats {
	stats := CorrectionStats{}
	for _, correction := range corrections {
		if correction.Corrected {
			stats.Corrected = append(stats.Corrected, correction)
		} else {
			stats.Pending = append(stats.Pending, correction)
		}
	}

	return stats
}

func FilterCorrections(values []ilias.Correction, test func(correction ilias.Correction) bool) (ret []ilias.Correction) {
	for _, s := range values {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
