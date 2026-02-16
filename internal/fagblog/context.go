package fagblog

import "text/template"

type Context struct {
	SiteMetadata SiteMetadata
	Templates map[string]*template.Template
}
