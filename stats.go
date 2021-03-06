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
	publishedPosts := GetPublishedPosts()
	draftPosts := GetDraftPosts()

	fmt.Printf("Number of published posts: %d\nNumber of drafts: %d\n", len(publishedPosts), len(draftPosts))
}

// ListDrafts will list the title of all Posts that are draft = true
func ListDrafts() {
	posts := GetDraftPosts()
	for _, p := range posts {
		fmt.Println(p.Title())
	}
}

// GetDraftPosts is a helper function purely for readability
// It issues a request to the getPosts function with the draft
// flage set to TRUE
func GetDraftPosts() []Post {
	return getPosts(true)
}

// GetPublishedPosts is a helper function purely for readability
// It issues a request to the getPosts function with the draft
// flag set to FALSE
func GetPublishedPosts() []Post {
	return getPosts(false)
}

// getPosts return an array of Post's that are in the proper draft state
// ordered by their date DESCENDING
func getPosts(isDraft bool) []Post {
	var posts []Post
	fileInfos, err := ioutil.ReadDir(MarkdownPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			post, err := NewPostFromFile(info)

			if err != nil {
				log.Fatal(err)
			}

			if post.Draft() == isDraft {
				posts = append(posts, post)
			}
		}
	}

	sort.Sort(PostsByDate(posts))
	return posts
}
