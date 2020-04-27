package util

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const (
	WorkspaceFilename = ".workspace.yml"
)

type WorkSpace struct {
	Exercise string
	Assignment string
	Corrections map[string][]string
}

func GetWorkspace() WorkSpace {
	data, err := ioutil.ReadFile(WorkspaceFilename)
	if err != nil {
		log.Fatal(err)
	}

	workspace := WorkSpace{}
	err = yaml.Unmarshal(data, &workspace)
	if err != nil {
		log.Fatal(err)
	}

	return workspace
}