package internal

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/JoeReid/openapi-route-optimiser/internal/openapi"
	"github.com/JoeReid/openapi-route-optimiser/internal/optimiser"
	"github.com/JoeReid/openapi-route-optimiser/internal/template"
	"github.com/go-playground/validator/v10"
	"github.com/jessevdk/go-flags"
)

// App holds the configuration and io writers for the application.
type App struct {
	SpecFile string `validate:"required,filepath" short:"s" long:"spec" description:"Path to the openapi file"`
	Debug    bool   `short:"d" long:"debug" description:"Should the program emit debug logs on stderr"`
	Filter   string `long:"filter" description:"Filter operationId tags using a regular expression. Executes before find and replace actions occour"`
	Find     string `long:"find" description:"Find sub-strings by regular expression to select for replacement"`
	Replace  string `long:"replace" description:"Replace found sub-strings with the given string. Supports capture groups from the regex"`
	Template string `validate:"omitempty,filepath" short:"t" long:"template" description:"Path to a go template file to format the output"`

	StdOut io.Writer
	StdErr io.Writer
}

// Run executes the main logic of the App.
func (a App) Run() {
	spec, err := openapi.LoadSpec(a.SpecFile)
	if err != nil {
		a.errorf("failed to load openapi spec file, %s", err.Error())
		os.Exit(1)
	}

	router := optimiser.NewRouter()

	a.debugf("processing the openapi spec...")
	for path, pathSpec := range spec.Paths {
		bindings := make(map[string]string)

		for verb, operation := range pathSpec.Operations() {
			a.debugf("processing %s %s %q", verb, path, operation.OperationId)

			opTags, err := processTags(operation.Tags, a.Filter, a.Find, a.Replace)
			if err != nil {
				a.errorf(err.Error())
				os.Exit(1)
			}

			switch len(opTags) {
			case 0:
				a.errorf("could not determine binding for operation %q (%s %s)", operation.OperationId, verb, path)
				os.Exit(1)
			case 1:
				a.debugf("found binding %q for operation %q", opTags[0], operation.OperationId)
				bindings[verb] = opTags[0]
			default:
				a.errorf("found multiple bindings (%q) for operation %q (%s %s)", strings.Join(opTags, ","), operation.OperationId, verb, path)
				os.Exit(1)
			}
		}

		// Make sure the bindings for all the verbs on a path are the same
		// TODO: should we support bindings per-verb?
		var verb, tag string
		for nextVerb, nextTag := range bindings {
			if verb != "" && nextTag != tag {
				a.errorf("found multiple bindings for path (%q): %s -> %s, %s -> %s", path, verb, tag, nextVerb, nextTag)
				os.Exit(1)
			}

			verb, tag = nextVerb, nextTag
		}

		// Strip the variables out of the paths, and replace with globbing patterns
		processedPath, err := processPath(path)
		if err != nil {
			a.errorf(err.Error())
			a.errorf("A serious internal error has occoured. Please report this issue!")
			os.Exit(1)
		}
		a.debugf("rewriting path %q -> %q", path, processedPath)

		a.debugf("adding path %q to route optimiser", processedPath)
		router.Subrouter(processedPath).Bind(tag)
	}

	a.debugf("optimising routes...")
	router.Prune()

	a.debugf("building the template payload...")
	var payload template.Payload
	router.Walk(func(path, service string) {
		payload = append(payload, template.Route{Path: path, Service: service})
	})

	a.debugf("executing template...")
	if err := template.Execute(a.StdOut, "", payload); err != nil {
		a.errorf(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func (a App) debugf(format string, args ...interface{}) {
	if a.Debug {
		fmt.Fprintf(a.StdErr, "DEBUG: "+format+"\n", args...)
	}
}

func (a App) errorf(format string, args ...interface{}) {
	fmt.Fprintf(a.StdErr, "ERROR: "+format+"\n", args...)
}

// New creates a new App instance with default values, over-riden by any runtime flags.
//
// New returns an error if we fail to process runtime flags, or if validation of the config fails.
func New() (*App, error) {
	app := &App{
		SpecFile: "openapi.yaml",
		Debug:    false,
		Filter:   "",
		Find:     "",
		Replace:  "",
		StdOut:   os.Stdout,
		StdErr:   os.Stderr,
	}

	// Read in the config from runtime flags
	if _, err := flags.Parse(app); err != nil {
		return nil, err
	}

	// Validate the config
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(app); err != nil {
		return nil, err
	}

	return app, nil
}

func processTags(tags []string, filter, find, replace string) ([]string, error) {
	filteredTags := make([]string, len(tags))
	copy(filteredTags, tags)

	if filter != "" {
		r, err := regexp.Compile(filter)
		if err != nil {
			return nil, fmt.Errorf("failed to compile filter regex, %s", err.Error())
		}

		// Filter the tags in-place
		n := 0
		for _, tag := range filteredTags {
			if r.MatchString(tag) {
				filteredTags[n] = tag
				n++
			}
		}
		filteredTags = filteredTags[:n]
	}

	if find != "" {
		r, err := regexp.Compile(find)
		if err != nil {
			return nil, fmt.Errorf("failed to compile find regex, %s", err.Error())
		}

		for i := range filteredTags {
			filteredTags[i] = r.ReplaceAllString(filteredTags[i], replace)
		}
	}

	return filteredTags, nil
}

func processPath(path string) (string, error) {
	r, err := regexp.Compile("{([a-zA-Z0-9]+_*-*)+}")
	if err != nil {
		return "", fmt.Errorf("failed to compile internal regex, %s", err.Error())
	}

	return r.ReplaceAllString(path, "*"), nil
}
