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

type Flags struct {
	templatePath string
	dataPath     string
	templateKey  string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flags := parseFlags()

	data, err := loadData(flags.dataPath)
	if err != nil {
		return fmt.Errorf("data loading error: %w", err)
	}

	tmpl, err := loadTemplate(flags.templatePath, flags.templateKey)
	if err != nil {
		return fmt.Errorf("template loading error: %w", err)
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}

	return nil
}

func parseFlags() Flags {
	templatePath := flag.String("t", "template.tmpl", "Path to the Go text template file")
	dataPath := flag.String("d", "data.yaml", "Path to the JSON or YAML file containing data")
	templateKey := flag.String("k", "", "Read template from structured file using specified key")
	flag.Parse()

	return Flags{
		templatePath: *templatePath,
		dataPath:     *dataPath,
		templateKey:  *templateKey,
	}
}

func loadTemplate(templatePath string, templateKey string) (*template.Template, error) {
	if templateKey == "" {
		return template.ParseFiles(templatePath)
	}

	templateFile, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	templateStr, err := parseTemplateFile(templatePath, templateFile, templateKey)
	if err != nil {
		return nil, err
	}

	return template.New("template").Parse(templateStr)
}

func parseTemplateFile(
	templatePath string,
	templateFile []byte,
	templateKey string,
) (string, error) {
	var templateData map[string]interface{}

	switch {
	case strings.HasSuffix(templatePath, ".yaml"), strings.HasSuffix(templatePath, ".yml"):
		if err := yaml.Unmarshal(templateFile, &templateData); err != nil {
			return "", fmt.Errorf("failed to parse template YAML: %w", err)
		}
	case strings.HasSuffix(templatePath, ".json"):
		if err := json.Unmarshal(templateFile, &templateData); err != nil {
			return "", fmt.Errorf("failed to parse template JSON: %w", err)
		}
	default:
		return "", errors.New("unsupported template file type")
	}

	templateStr, ok := templateData[templateKey].(string)
	if !ok {
		return "", fmt.Errorf("template key '%s' not found or not a string", templateKey)
	}

	return templateStr, nil
}

func loadData(dataPath string) (interface{}, error) {
	dataFile, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var data interface{}
	switch {
	case strings.HasSuffix(dataPath, ".yaml"), strings.HasSuffix(dataPath, ".yml"):
		if err = yaml.Unmarshal(dataFile, &data); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	case strings.HasSuffix(dataPath, ".json"):
		if err = json.Unmarshal(dataFile, &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return nil, errors.New("unsupported data file type")
	}

	return data, nil
}
