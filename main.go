package main

import (
	"fmt"
	"os"

	"github.com/mrnickel/StaticSiteGenerator/publish"
	"github.com/mrnickel/StaticSiteGenerator/stats"
)

func main() {

	command := os.Args[1]

	switch command {
	case "publish":
		fmt.Println("Publish the markdown file specified")
		publish.Publish(os.Args[2])
		return
	case "create":
		fmt.Println("Create a Post")
		CreatePost(os.Args[2])
		return
	case "stats":
		fmt.Println("Get the stats for this site")
		stats.GetStats()
		return
	case "listdrafts":
		stats.ListDrafts()
		return
	}
}
