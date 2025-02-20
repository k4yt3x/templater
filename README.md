# Templater

A simple command-line Go `text/template` rendering utility.

This tool provides an easy way to write and troubleshoot Go `text/template` files. It reads a standard Go `text/template` file and a YAML or JSON file containing the data to render the template with, then writes the rendered template to standard output.

I know it's extremely simple, but I still find it useful.

![screenshot](https://github.com/user-attachments/assets/33015d22-5de8-492e-83dd-829d91da6035)

## Usage

```bash
templater -t template.tmpl -d data.yaml
```

## Build

```bash
go build -o bin/templater ./cmd/templater/...
```
