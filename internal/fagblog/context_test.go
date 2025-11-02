package fagblog

import (
	"sync"
	"testing"
)

func TestPostCache_ThreadSafety(t *testing.T) {
	cache := NewPostCache()
	
	post := BlogPost{
		Metadata: BlogPostMetadata{
			Title: "Test Post",
		},
	}
	
	metadata := BlogPostMetadata{
		Title: "Test Metadata",
	}
	
	// Test concurrent writes and reads
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		
		// Concurrent writes
		go func(i int) {
			defer wg.Done()
			cache.SetPost("test-post", post)
			cache.SetMetadata("test-metadata", metadata)
		}(i)
		
		// Concurrent reads
		go func(i int) {
			defer wg.Done()
			cache.GetPost("test-post")
			cache.GetMetadata("test-metadata")
		}(i)
	}
	
	wg.Wait()
	
	// Verify final state
	cachedPost, ok := cache.GetPost("test-post")
	if !ok {
		t.Error("Expected post to be in cache")
	}
	if cachedPost.Metadata.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", cachedPost.Metadata.Title)
	}
	
	cachedMetadata, ok := cache.GetMetadata("test-metadata")
	if !ok {
		t.Error("Expected metadata to be in cache")
	}
	if cachedMetadata.Title != "Test Metadata" {
		t.Errorf("Expected title 'Test Metadata', got '%s'", cachedMetadata.Title)
	}
}

func TestPostCache_GetSetPost(t *testing.T) {
	cache := NewPostCache()
	
	// Test cache miss
	_, ok := cache.GetPost("non-existent")
	if ok {
		t.Error("Expected cache miss for non-existent post")
	}
	
	// Test cache hit
	post := BlogPost{
		Metadata: BlogPostMetadata{
			Title: "Test Post",
		},
	}
	cache.SetPost("test", post)
	
	cached, ok := cache.GetPost("test")
	if !ok {
		t.Error("Expected cache hit for existing post")
	}
	if cached.Metadata.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", cached.Metadata.Title)
	}
}

func TestPostCache_GetSetMetadata(t *testing.T) {
	cache := NewPostCache()
	
	// Test cache miss
	_, ok := cache.GetMetadata("non-existent")
	if ok {
		t.Error("Expected cache miss for non-existent metadata")
	}
	
	// Test cache hit
	metadata := BlogPostMetadata{
		Title: "Test Metadata",
	}
	cache.SetMetadata("test", metadata)
	
	cached, ok := cache.GetMetadata("test")
	if !ok {
		t.Error("Expected cache hit for existing metadata")
	}
	if cached.Title != "Test Metadata" {
		t.Errorf("Expected title 'Test Metadata', got '%s'", cached.Title)
	}
}
