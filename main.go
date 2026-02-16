package main

import (
	"fmt"
	"text/template"
	"log"
	"net/http"
	"os"

	"github.com/airkoala/fagblog/internal/fagblog"
	"github.com/airkoala/fagblog/internal/middleware"
	"github.com/airkoala/fagblog/internal/route"
)

func main() {
	config, err := fagblog.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	blogMetadata, err := fagblog.SiteMetadataFromToml(config.ContentDir + "/meta.toml")

	if err != nil {
		log.Fatalf("Error loading blog metadata: %v", err)
	}

	mux := http.NewServeMux()

	context := fagblog.Context{
		SiteMetadata: blogMetadata,
		Templates:    loadTemplates(config.TemplateDir),
	}

	handle(mux, route.Home(&context, &config))
	handle(mux, route.Static(&context, &config))
	handle(mux, route.Assets(&context, &config))
	handle(mux, route.BlogPost(&context, &config))
	handle(mux, route.Feed(&context, &config))

	log.Printf("Starting server on :%d", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux)

	log.Fatalln("Server stopped: ", err)
}

func handle(s *http.ServeMux, route route.Route) {
	middlewares := route.Middlewares

	// Add default middlewares
	middlewares = append(middlewares, middleware.Logging())

	log.Printf("Registering %s\n", route.Pattern)

	// Register the route chained with middlewares
	s.HandleFunc(route.Pattern, middleware.Chain(route.Handler, middlewares...))
}

// Loads templates from the specified directory and associates them with their layouts.
func loadTemplates(templateDir string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	// Base template with all layouts associated.
	layout := template.Must(template.ParseGlob(templateDir + "/layout/*"))

	extendedTemplates, err := os.ReadDir(templateDir + "/extends/")
	if err != nil {
		log.Fatalf("Error reading template directory: %v\n", err)
	}

	// Iterate over all template files in <templateDir>/extends/*, skipping directories.
	for _, e := range extendedTemplates {
		if e.IsDir() {
			continue
		}

		templatePath := templateDir + "/extends/" + e.Name()

		// .Clone() shallow copies the template without duplicating associated (layout) templates.
		tmpl, err := template.Must(layout.Clone()).ParseFiles(templatePath)
		if err != nil {
			log.Fatalf("Error parsing template %s: %v\n", templatePath, err)
		}

		templates[e.Name()] = tmpl
		log.Printf("Loaded template: %s\n", e.Name())
	}

	// Templates which do not extend the layout
	baseTemplates, err := os.ReadDir(templateDir + "/base/")
	if err != nil {
		log.Fatalf("Error reading template directory: %v\n", err)
	}

	for _, e := range baseTemplates {
		templatePath := templateDir + "/base/" + e.Name()

		// .Clone() shallow copies the template without duplicating associated (layout) templates.
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			log.Fatalf("Error parsing template %s: %v\n", templatePath, err)
		}

		templates[e.Name()] = tmpl
		log.Printf("Loaded template: %s\n", e.Name())
	}


	return templates
}
