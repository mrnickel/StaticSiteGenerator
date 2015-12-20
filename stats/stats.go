package stats

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/mrnickel/StaticSiteGenerator/constants"
	"github.com/mrnickel/StaticSiteGenerator/post"
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
		fmt.Println(p.Title)
	}
}

// GetDraftPosts is a helper function purely for readability
// It issues a request to the getPosts function with the draft
// flage set to TRUE
func GetDraftPosts() []*post.Post {
	return getPosts(true)
}

// GetPublishedPosts is a helper function purely for readability
// It issues a request to the getPosts function with the draft
// flag set to FALSE
func GetPublishedPosts() []*post.Post {
	return getPosts(false)
}

// getPosts return an array of Post's that are in the proper draft state
// ordered by their date DESCENDING
func getPosts(isDraft bool) []*post.Post {
	var posts []*post.Post
	fileInfos, err := ioutil.ReadDir(constants.MarkdownPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			post := post.NewPostFromFile(info)
			if post.Draft == isDraft {
				posts = append(posts, post)
			}
		}
	}

	sort.Sort(post.PostsByDate(posts))
	return posts
}
