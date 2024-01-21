package template

import (
	"io"
	"text/template"
)

func Execute(w io.Writer, file string, data Payload) error {
	if file == "" {
		return defaultTemplate.Execute(w, data)
	}

	t, err := template.ParseFiles(file)
	if err != nil {
		return err
	}

	return t.Execute(w, data)
}

type Payload []Route

func (p Payload) Services() map[string][]string {
	rtn := make(map[string][]string)

	for _, route := range p {
		if paths, ok := rtn[route.Service]; ok {
			rtn[route.Service] = append(paths, route.Path)
			continue
		}

		rtn[route.Service] = []string{route.Path}
	}

	return rtn
}

func (p Payload) Paths() map[string]string {
	rtn := make(map[string]string)

	for _, route := range p {
		rtn[route.Path] = route.Service
	}

	return rtn
}

type Route struct {
	Path    string
	Service string
}

var defaultTemplate = template.Must(template.New("default").Parse(`
{{- range $path, $service := .Paths}}
	{{- $path}} -> {{$service}}
{{end -}}
`))
