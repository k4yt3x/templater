package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Define flags for data file and template file paths
	templatePath := flag.String("t", "template.tmpl", "Path to the Go text template file")
	dataPath := flag.String("d", "data.yaml", "Path to the JSON or YAML file containing data")
	flag.Parse()

	// Read data file
	dataFile, err := os.ReadFile(*dataPath)
	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}

	// Determine file type and parse accordingly
	var data map[string]interface{}
	switch {
	case strings.HasSuffix(*dataPath, ".yaml"), strings.HasSuffix(*dataPath, ".yml"):
		if err = yaml.Unmarshal(dataFile, &data); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
	case strings.HasSuffix(*dataPath, ".json"):
		if err = json.Unmarshal(dataFile, &data); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return errors.New("unsupported file type")
	}

	// Read and parse template file
	tmpl, err := template.ParseFiles(*templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with parsed data
	if err = tmpl.Execute(os.Stdout, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
