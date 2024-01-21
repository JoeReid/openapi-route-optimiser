package openapi

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Spec struct {
	Paths map[string]Path `yaml:"paths"`
}

type Path struct {
	Get     *Operation `yaml:"get"`
	Post    *Operation `yaml:"post"`
	Put     *Operation `yaml:"put"`
	Patch   *Operation `yaml:"patch"`
	Delete  *Operation `yaml:"delete"`
	Head    *Operation `yaml:"head"`
	Options *Operation `yaml:"options"`
}

func (p Path) Operations() map[string]Operation {
	rtn := make(map[string]Operation)

	if p.Get != nil {
		rtn["GET"] = *p.Get
	}
	if p.Post != nil {
		rtn["POST"] = *p.Post
	}
	if p.Put != nil {
		rtn["PUT"] = *p.Put
	}
	if p.Patch != nil {
		rtn["PATCH"] = *p.Patch
	}
	if p.Delete != nil {
		rtn["DELETE"] = *p.Delete
	}
	if p.Head != nil {
		rtn["HEAD"] = *p.Head
	}
	if p.Options != nil {
		rtn["OPTIONS"] = *p.Options
	}

	return rtn
}

type Operation struct {
	OperationId string   `yaml:"operationId"`
	Tags        []string `yaml:"tags"`
}

func LoadSpec(filename string) (*Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var spec Spec
	if err := yaml.Unmarshal(bytes, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}
