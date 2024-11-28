package main

import (
	"flag"
	"html"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type SectionItem struct {
	Term        string `yaml:"term"`
	Description string `yaml:"description"`
}

type Section struct {
	Title string        `yaml:"title"`
	Items []SectionItem `yaml:"items"`
}

type PageData struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Sections    []Section `yaml:"sections"`
}

func ParsePageData(path string) PageData {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed reading %s: %v", path, err)
	}

	var pageData PageData
	err = yaml.Unmarshal(data, &pageData)
	if err != nil {
		log.Fatalf("Failed parsing %s: %v", path, err)
	}
	return encodePageData(pageData)
}

func encodePageData(data PageData) PageData {
	data.Title = html.EscapeString(data.Title)
	data.Description = html.EscapeString(data.Description)
	for i := 0; i < len(data.Sections); i++ {
		data.Sections[i].Title = html.EscapeString((data.Sections[i].Title))
		for j := 0; j < len(data.Sections[i].Items); j++ {
			data.Sections[i].Items[j].Term = html.EscapeString(data.Sections[i].Items[j].Term)
			data.Sections[i].Items[j].Description = html.EscapeString(data.Sections[i].Items[j].Description)
		}
	}
	return data
}

func main() {
	var inputFile string
	var outputFile string
	var templateName string

	flag.StringVar(&inputFile, "i", "", "Input file")
	flag.StringVar(&inputFile, "input", "", "Input file")

	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&outputFile, "output", "", "Output file")

	flag.StringVar(&templateName, "t", "", "Template name")
	flag.StringVar(&templateName, "template", "printable", "Template name")

	flag.Parse()

	if inputFile == "" {
		log.Fatal("Input file is required")
	}

	if !(strings.HasSuffix(inputFile, ".yaml") || strings.HasSuffix(inputFile, ".yml")) {
		log.Fatal("Input file must be a YAML file")
	}

	if outputFile == "" {
		log.Fatal("Output file is required")
	}

	generateCheatsheet(inputFile, outputFile, templateName)
}

func generateCheatsheet(inputFile string, outputFile string, templateName string) {
	templ := template.Must(template.ParseFiles("templates/" + templateName + ".html"))

	if !strings.HasSuffix(outputFile, ".html") {
		outputFile += ".html"
	}

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed creating output file: %v", err)
	}

	defer file.Close()

	err = templ.Execute(file, ParsePageData(inputFile))
	if err != nil {
		log.Fatalf("Failed executing the template: %v", err)
	}
}
