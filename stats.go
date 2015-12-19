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
	publishedArticles := getPublishedArticles()
	draftArticles := getDraftArticles()

	fmt.Printf("Number of published posts: %d\nNumber of drafts: %d\n", len(publishedArticles), len(draftArticles))
}

// ListDrafts will list the title of all posts that are draft = true
func ListDrafts() {
	articles := getDraftArticles()
	for _, article := range articles {
		fmt.Println(article.Title)
	}
}

// getDraftArticles is a helper function purely for readability
// It issues a request to the getArticles function with the draft
// flage set to TRUE
func getDraftArticles() []*Post {
	return getArticles(true)
}

// getPublishedArticles is a helper function purely for readability
// It issues a request to the getArticles function with the draft
// flag set to FALSE
func getPublishedArticles() []*Post {
	return getArticles(false)
}

// getArticles return an array of Post's that are in the proper draft state
// ordered by their date DESCENDING
func getArticles(isDraft bool) []*Post {
	var articles []*Post
	fileInfos, err := ioutil.ReadDir(baseArticlePath)

	if err != nil {
		log.Fatal(err)
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			post := NewPostFromFile(info)
			if post.Draft == isDraft {
				articles = append(articles, NewPostFromFile(info))
			}
		}
	}

	sort.Sort(PostsByDate(articles))
	return articles
}
