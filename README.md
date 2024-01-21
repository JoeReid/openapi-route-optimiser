[![Go](https://github.com/JoeReid/openapi-route-optimiser/actions/workflows/go.yml/badge.svg)](https://github.com/JoeReid/openapi-route-optimiser/actions/workflows/go.yml)

# OpenAPI Route Optimiser

This project is a tool for extracting and optimising routing rules for an OpenAPI yaml specification.

## Usage

```
$ ./openapi-route-optimiser --help

Usage:
  main [OPTIONS]

Application Options:
  -s, --spec=     Path to the openapi file (default: openapi.yaml)
  -d, --debug     Should the program emit debug logs on stderr
      --filter=   Filter operationId tags using a regular expression. Executes before find and replace actions occour
      --find=     Find sub-strings by regular expression to select for replacement
      --replace=  Replace found sub-strings with the given string. Supports capture groups from the regex
  -t, --template= Path to a go template file to format the output

Help Options:
  -h, --help      Show this help message

```

## Installation

### From source

```
git clone https://github.com/JoeReid/openapi-route-optimiser.git
cd openapi-route-optimiser
go build
```

This will create a openapi-route-optimiser binary in the current directory.
