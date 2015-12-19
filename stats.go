package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// GetStats will print out a list of posts in the various states
// namely how many are in draft state (draft = true), and
// how many are in published state (draft = false)
func GetStats() {
	publishedPosts := getPublishedPosts()
	draftPosts := getDraftPosts()

	fmt.Printf("Number of published posts: %d\nNumber of drafts: %d\n", len(publishedPosts), len(draftPosts))
}

// ListDrafts will list the title of all Posts that are draft = true
func ListDrafts() {
	posts := getDraftPosts()
	for _, p := range posts {
		fmt.Println(p.Title)
	}
}

// getDraftPosts is a helper function purely for readability
// It issues a request to the getPosts function with the draft
// flage set to TRUE
func getDraftPosts() []*Post {
	return getPosts(true)
}

// getPublishedPosts is a helper function purely for readability
// It issues a request to the getPosts function with the draft
// flag set to FALSE
func getPublishedPosts() []*Post {
	return getPosts(false)
}

// getPosts return an array of Post's that are in the proper draft state
// ordered by their date DESCENDING
func getPosts(isDraft bool) []*Post {
	var posts []*Post
	fileInfos, err := ioutil.ReadDir(baseMarkdownPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			post := NewPostFromFile(info)
			if post.Draft == isDraft {
				posts = append(posts, NewPostFromFile(info))
			}
		}
	}

	sort.Sort(PostsByDate(posts))
	return posts
}
