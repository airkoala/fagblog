package fagblog

import (
	"html/template"
	"sync"
)

type Context struct {
	SiteMetadata SiteMetadata
	Templates    map[string]*template.Template
	PostCache    *PostCache
}

type PostCache struct {
	mu       sync.RWMutex
	posts    map[string]BlogPost
	metadata map[string]BlogPostMetadata
}

func NewPostCache() *PostCache {
	return &PostCache{
		posts:    make(map[string]BlogPost),
		metadata: make(map[string]BlogPostMetadata),
	}
}

func (pc *PostCache) GetPost(key string) (BlogPost, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	post, ok := pc.posts[key]
	return post, ok
}

func (pc *PostCache) SetPost(key string, post BlogPost) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.posts[key] = post
}

func (pc *PostCache) GetMetadata(key string) (BlogPostMetadata, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	metadata, ok := pc.metadata[key]
	return metadata, ok
}

func (pc *PostCache) SetMetadata(key string, metadata BlogPostMetadata) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.metadata[key] = metadata
}
