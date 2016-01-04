package main

import (
	"fmt"
	"os"

	"github.com/mrnickel/StaticSiteGenerator/constants"
	"github.com/mrnickel/StaticSiteGenerator/post"
)

// CreatePost creates a blog post
func CreatePost(title string) {
	file, err := os.Create(constants.MarkdownPath + fmt.Sprintf("%s.md", post.GenerateFileNamePrefix(title)))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := post.NewPost(title)

	file.WriteString(p.String())
}
