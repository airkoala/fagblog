package fagblog

import (
	"bytes"
	"errors"
	"text/template"
	"log"
	"os"
	"slices"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/net/html"
)

type Heading struct {
	Title string
	Level uint
	Id    string
}

type BlogPostMetadata struct {
	Title         string
	Timestamp     time.Time
	Summary       string
	ThumbnailHref string
}

type BlogPost struct {
	Metadata BlogPostMetadata
	Content  string
	Headings []Heading
}

// Parses a TOML file and returns a BlogPost struct.
func GetPost(dirPath string, postName string) (BlogPost, error) {
	post := BlogPost{}
	postDirPath := dirPath + "/" + postName

	// Check if the post directory exists
	// If it doesn't exist, return os.ErrNotExist
	if _, err := os.Stat(postDirPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("Post directory does not exist: %s", postDirPath)
		return post, err
	}

	metadata, err := GetPostMetadata(dirPath, postName)
	if err != nil {
		return post, err
	}

	post.Metadata = metadata

	// content, err := os.ReadFile(postDirPath + "/index.html")
	// if err != nil {
	// 	log.Printf("Error reading file %s: %v", postDirPath+"/index.html", err)
	// }

	content, headings, err := getPostContentAndHeadings(postDirPath + "/index.html")
	if err != nil {
		log.Printf("Error reading file %s: %v", postDirPath+"/index.html", err)
		return post, err
	}

	post.Content = content
	post.Headings = headings

	// // Debug
	// for _, h := range headings {
	// 	log.Printf("Heading: %s, Level: %d, Id: %s", h.Title, h.Level, h.Id)
	// }

	return post, nil
}

func GetPostMetadata(dirPath string, postName string) (BlogPostMetadata, error) {
	metadata := BlogPostMetadata{}
	postDirPath := dirPath + "/" + postName

	// Check if the post directory exists
	// If it doesn't exist, return os.ErrNotExist
	if _, err := os.Stat(postDirPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("Post directory does not exist: %s", postDirPath)
		return metadata, err
	}

	_, err := toml.DecodeFile(postDirPath+"/meta.toml", &metadata)

	if err != nil {
		log.Printf("Error decoding TOML file: %v", err)
		return metadata, err
	}

	return metadata, nil
}

// GetPosts returns a list of all posts in the specified directory.
func GetPosts(dirPath string) ([]string, error) {
	posts := make([]string, 0)

	// Check if the post directory exists
	// If it doesn't exist, return os.ErrNotExist
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("Post directory does not exist: %s", dirPath)
		return posts, err
	}

	entries, err := os.ReadDir(dirPath)

	if err != nil {
		log.Printf("Error reading directory %s: %v", dirPath, err)
		return posts, err
	}

	// to ensure that posts are in order of newest to oldest
	slices.Reverse(entries)

	for _, entry := range entries {
		if entry.IsDir() {
			posts = append(posts, entry.Name())
		}
	}

	return posts, nil
}

// Gets the parsed and formatted post content with headings.
func getPostContentAndHeadings(path string) (string, []Heading, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v", path, err)
		return "", []Heading{}, err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		log.Printf("Error parsing HTML file %s: %v", path, err)
		return "", []Heading{}, err
	}

	// This also inserts ids for hX tags.
	// Yes its janky deal with it.
	headings, err := createToc(doc)
	if err != nil {
		log.Printf("Error creating table of contents: %v", err)
		return "", []Heading{}, err
	}

	var buf bytes.Buffer
	html.Render(&buf, doc)

	return buf.String(), headings, nil
}

func createToc(doc *html.Node) ([]Heading, error) {
	headings := make([]Heading, 0)
	headingTags := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	hNodes := make([]*html.Node, 0)

	// Recursively walk child nodes and collect headings.
	for n := range doc.Descendants() {
		// Skip non element nodes.
		if n.Type != html.ElementNode {
			continue
		}

		if slices.Contains(headingTags, n.Data) {
			hNodes = append(hNodes, n)
		}
	}

	for _, n := range hNodes {
		id, err := getOrCreateId(n)
		if err != nil {
			log.Printf("Error getting or creating id: %v", err)
			continue
		}

		headings = append(headings, Heading{
			Title: n.FirstChild.Data, // FirstChild is guaranteed to be a TextNode
			Id:    id,

			// Cast is safe since tag is guaranteed to be in array.
			Level: uint(slices.Index(headingTags, n.Data)),
		})
	}

	normaliseHeadings(headings)

	return headings, nil
}

// Returns the id of a node or creates a url encoded id from the node's content and adds it to the node.
func getOrCreateId(node *html.Node) (string, error) {
	for _, attr := range node.Attr {
		// Check if an id is already present or not
		if attr.Key == "id" {
			return attr.Val, nil
		}
	}

	if node.FirstChild.Type != html.TextNode {
		return "", errors.New("First child is not a text node")
	}

	id := template.URLQueryEscaper(node.FirstChild.Data)
	node.Attr = append(node.Attr, html.Attribute{Key: "id", Val: id})

	return id, nil
}

// Normalises heading levels to be 0, 1, 2, etc.
func normaliseHeadings(headings []Heading) {
	// Get a set of levels
	uniques := make(map[uint]bool)
	for _, h := range headings {
		uniques[h.Level] = true
	}

	// Store levels into a sorted slice
	levels := make([]uint, 0, len(uniques))
	for k := range uniques {
		levels = append(levels, k)
	}
	slices.Sort(levels)

	// Set heading levels to normalised values.
	for _, h := range headings {
		h.Level = uint(slices.Index(levels, h.Level))
	}
}
