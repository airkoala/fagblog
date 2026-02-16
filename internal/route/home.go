package route

import (
	"log"
	"net/http"

	"github.com/airkoala/fagblog/internal/fagblog"
)

type post struct {
		Id      string
		Metadata fagblog.BlogPostMetadata
	}
type homeData struct {
	Context *fagblog.Context
	Posts   []post 
	Url string
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

			homeData := homeData{
				Context: context,
				Posts: make([]post, 0, len(postNames)),
				Url: r.URL.String(),
			}

			for _, n := range postNames {
				metadata, err := fagblog.GetPostMetadata(config.ContentDir+"/blog", n)
				if err != nil {
					log.Printf("Error getting post metadata: %v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				homeData.Posts = append(homeData.Posts, post{
					Id: n,
					Metadata: metadata,
				})
			}

			err = context.Templates["home.html"].Execute(w, homeData)

			if err != nil {
				log.Printf("Error executing template: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}}
}
