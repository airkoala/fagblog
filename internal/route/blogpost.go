package route

import (
	"errors"
	// "html/template"
	"log"
	"net/http"
	"os"

	"github.com/airkoala/fagblog/internal/fagblog"
)

type BlogData struct {
	Context *fagblog.Context
	Post    fagblog.BlogPost
	Url     string
}

func BlogPost(context *fagblog.Context, config *fagblog.Config) Route {
	return Route{
		Pattern: "GET /blog/{postName}",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			postName := r.PathValue("postName")
			post, err := fagblog.GetPost(config.ContentDir+"/blog", postName, context.PostCache)

			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("Post not found: %v\n", err)
					http.Error(w, "Post not found", http.StatusNotFound)
					return
				}
				log.Printf("Error getting post: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// tmpl := template.Must(template.New("blogpost.html").ParseFiles("templates/blogpost.html", "templates/base.html", "templates/header.html", "templates/footer.html"))

			err = context.Templates["blogpost.html"].Execute(w, BlogData{
				Context: context,
				Post:    post,
				Url:     r.URL.String(),
			})
			// err = context.Templates.ExecuteTemplate(w, "blogpost.html", BlogData{Context: context, Post: post})

			if err != nil {
				log.Printf("Error executing template: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		},
	}
}
