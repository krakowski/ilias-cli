package util

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const (
	WorkspaceFilename = ".workspace.yml"
)

type Exercise struct {
	Reference	string 		`yaml:"reference"`
	Assignment	string 		`yaml:"assignment"`
}

type Table struct {
	Reference	string 		`yaml:"reference"`
	Identifier	string 		`yaml:"identifier"`
	Name	    string		`yaml:"name"`
}

type WorkSpace struct {
	Exercise 		  Exercise					`yaml:"exercise"`
	Table 		  	  Table						`yaml:"table"`
	Corrections       map[string][]string		`yaml:"corrections"`
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