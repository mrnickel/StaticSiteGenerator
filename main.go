package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	if len(os.Args) <= 1 || os.Args[1] == "help" {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "publish":
		tmpP := NewPost(os.Args[2])

		fmt.Println("Publish the markdown file specified")
		file, err := os.Open(tmpP.MarkdownPath())
		if err != nil {
			log.Fatal(err)
		}
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		p, err := NewPostFromFile(fileInfo)
		if err != nil {
			log.Fatal(err)
		}
		err = p.Publish()
		if err != nil {
			log.Fatal(err)
		}
		return
	case "create":
		fmt.Println("Create a Post")
		p := NewPost(os.Args[2])
		p.Update()
		return
	case "stats":
		fmt.Println("Get the stats for this site")
		GetStats()
		return
	case "listdrafts":
		ListDrafts()
		return
	case "newsite":
		fmt.Println("TODO!")
		return
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("You must choose one of the following options:\n")
	fmt.Println("publish \"Blog Title here\"")
	fmt.Println("create \"Blog Title here\"")
	fmt.Println("stats (this will list stats about your site)")
	fmt.Println("listdrafts (this will list the titles of all your posts still in draft mode)")
	fmt.Println("newsite (creates a new site -- still needs to be implemented)")
}
