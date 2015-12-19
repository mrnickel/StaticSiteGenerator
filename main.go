package main

import (
	"fmt"
	"os"
)

const baseArticlePath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/posts/"
const baseTemplatePath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/templates/"
const basePublushPath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/articles/"
const baseRootPath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/"

func main() {

	command := os.Args[1]

	switch command {
	case "publish":
		fmt.Println("Publish the markdown file specified")
		Publish(os.Args[2])
		return
	case "create":
		fmt.Println("Create a blog post")
		CreatePost(os.Args[2])
		return
	case "stats":
		fmt.Println("Get the stats for this site")
		GetStats()
		return
	case "listdrafts":
		ListDrafts()
		return
	}
}
