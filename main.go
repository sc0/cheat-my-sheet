package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
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
	return pageData
}

func main() {
	templ := template.Must(template.ParseFiles("templates/default.html"))

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := ParsePageData("data/test.yaml")

		err := templ.Execute(w, data)
		if err != nil {
			log.Fatalf("Failed executing the template: %v", err)
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	srv.ListenAndServe()
}
