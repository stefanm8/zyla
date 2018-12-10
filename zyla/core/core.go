package core

import (
	"github.com/asdine/storm"
	"github.com/gorilla/mux"
)

const DatabaseName = "bolt.db"

type Dependency struct {
	Import string
	Name   string
	User   string
}

type Field struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Options string `yaml:"options"`
}

type Model struct {
	BucketName string   `yaml:"bucket"`
	Fields     []*Field `yaml:"fields"`
	Identifier string   `yaml:"identifier"`
}

type Application struct {
	Name        string `yaml:"name"`
	Route       string `yaml:"route"`
	Model       Model  `yaml:"model"`
	Dependecies map[string]*Dependency
}

type Project struct {
	ProjectName  string         `yaml:"name"`
	Github       string         `yaml:"github"`
	Applications []*Application `yaml:"applications"`
	Dependecies  map[string]*Dependency
}

type Initializer struct {
	Project *Project
	DB      *storm.DB
	Router  *mux.Router
}
