package main

import (
	"fmt"
	"os"
	"strings"
)

// CreatePost creates a blog post
func CreatePost(title string) {
	file, err := os.Create(fmt.Sprintf("%s.md", strings.Replace(strings.ToLower(title), " ", "+", -1)))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := NewPost(title)

	file.WriteString(p.GetString())
}
