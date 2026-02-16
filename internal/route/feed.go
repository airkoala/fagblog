package route

import (
	"log"
	"net/http"
	// "time"

	"github.com/airkoala/fagblog/internal/fagblog"
)

type feedData struct {
	Context      *fagblog.Context
	Posts        map[string]fagblog.BlogPost
	Url          string

	// TODO:
	// LastUpdateTS string
}

const feedEntryContentDisclaimer = `<p> Warning: This page is best viewed on the browser.
Your feed reader may not render the content correctly.</p> <br/>`

func Feed(context *fagblog.Context, config *fagblog.Config) Route {
	return Route{
		Pattern: "GET /feed.atom",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			postNames, err := fagblog.GetPosts(config.ContentDir + "/blog")

			if err != nil {
				log.Printf("Error getting posts: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := feedData{
				Context: context,
				Posts:   make(map[string]fagblog.BlogPost, len(postNames)),
				Url:     r.URL.String(),
			}

			for _, n := range postNames {
				post, err := fagblog.GetPost(config.ContentDir+"/blog", n)
				if err != nil {
					log.Printf("Error getting post: %v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				post.Content = feedEntryContentDisclaimer + post.Content
				data.Posts[n] = post
			}

			w.Header().Set("Content-Type", "application/atom+xml")
			w.Header().Set("Cache-Control", "max-age=36000")  // 10 h cache
			err = context.Templates["feed.atom"].Execute(w, data)

			if err != nil {
				log.Printf("Error executing template: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}}
}
