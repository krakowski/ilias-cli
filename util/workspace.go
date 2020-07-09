package util

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

const (
	WorkspaceFilename = ".workspace.yml"
	UsercacheFilename = ".user"
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

type Workspace struct {
	Title			  string					`yaml:"title"`
	Exercise 		  Exercise					`yaml:"exercise"`
	//Table 		  	  Table						`yaml:"table"`
	Corrections       map[string][]string		`yaml:"corrections"`
}

func ReadWorkspace() Workspace {
	assertWorkspaceExists()
	data, err := ioutil.ReadFile(WorkspaceFilename)
	if err != nil {
		log.Fatal("reading workspace failed: " + err.Error())
	}

	workspace := Workspace{}
	err = yaml.Unmarshal(data, &workspace)
	if err != nil {
		log.Fatal("unmarshaling workspace faild: " + err.Error())
	}

	return workspace
}

func WriteUserCache(username string) {
	if err := ioutil.WriteFile(UsercacheFilename, []byte(username), 0644); err != nil {
		log.Fatal(err)
	}
}

func ReadUserCache() string {
	username, err := ioutil.ReadFile(UsercacheFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, Red("workspace is not initialized"))
		os.Exit(1)
	}

	return string(username)
}

func assertWorkspaceExists() {
	if _, err := os.Stat(WorkspaceFilename); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, Red("not a ILIAS workspace"))
		os.Exit(1)
	}
}