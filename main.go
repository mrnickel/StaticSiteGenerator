package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const baseArticlePath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/posts/"
const baseTemplatePath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/templates/"
const basePublushPath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/articles/"
const baseRootPath string = "/Users/ryannickel/Documents/Pending/mrnickel.github.io/"

// Article structs are very simple in that they only have content. i.e. what's in the markdown file.
type Article struct {
	Content     string
	PublishTime time.Time
}

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

	// tmpCheckoutPath, _ := ioutil.TempDir(os.TempDir(), "mrnickel")
	// articlePath := tmpCheckoutPath + "/posts/"

	// fetchGitRepo("https://github.com/mrnickel/mrnickel.github.io.git", tmpCheckoutPath)
	// articles := getArticles(articlePath)

	// for _, article := range articles {
	// 	post := parseArticle(articlePath + article.Name())
	// 	article := Article{
	// 		Content:     post,
	// 		PublishTime: time.Now(),
	// 	}

	// 	generateArticle(tmpCheckoutPath, article)

	// }

	// os.RemoveAll(tmpCheckoutPath)
}

// Uses github to check out the specified git repo
// Doesn't take into account any auth at the moment
func fetchGitRepo(url string, checkoutPath string) {
	cmd := exec.Command("git", "clone", url, ".")
	cmd.Dir = checkoutPath

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
