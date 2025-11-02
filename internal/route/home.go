package route

import (
	"github.com/airkoala/fagblog/internal/fagblog"
	"log"
	"net/http"
)

type HomeData struct {
	Context *fagblog.Context
	Posts   map[string]fagblog.BlogPostMetadata
	Url     string
}

func Home(context *fagblog.Context, config *fagblog.Config) Route {
	return Route{
		Pattern: "GET /{$}",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			postNames, err := fagblog.GetPosts(config.ContentDir + "/blog")

			if err != nil {
				log.Printf("Error getting posts: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			homeData := HomeData{
				Context: context,
				Posts:   make(map[string]fagblog.BlogPostMetadata, len(postNames)),
				Url:     r.URL.String(),
			}

			for _, n := range postNames {
				metadata, err := fagblog.GetPostMetadata(config.ContentDir+"/blog", n, context.PostCache)
				if err != nil {
					log.Printf("Error getting post metadata: %v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				homeData.Posts[n] = metadata
			}

			err = context.Templates["home.html"].Execute(w, homeData)

			if err != nil {
				log.Printf("Error executing template: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}}
}
