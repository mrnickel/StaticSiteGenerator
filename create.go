package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mrnickel/StaticSiteGenerator/config"
	"github.com/mrnickel/StaticSiteGenerator/post"
)

// CreatePost creates a blog post
func CreatePost(title string) {
	file, err := os.Create(config.MarkdownPath + fmt.Sprintf("%s.md", strings.Replace(strings.ToLower(title), " ", "+", -1)))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := post.NewPost(title)

	file.WriteString(p.String())
}
