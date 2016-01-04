package publish

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/mrnickel/StaticSiteGenerator/constants"
	"github.com/mrnickel/StaticSiteGenerator/post"
	"github.com/mrnickel/StaticSiteGenerator/stats"
)

// Publish will take the Post.Title of the item we wish to publish, generate the HTML from the
// template, update the .md file's draft flag to false and re-generate the index page.
// Maybe one day I'll add other flags, such as "+tweet" in order to connect to twitter
// and post on my behalf
func Publish(postTitle string) {
	mdFileName := fmt.Sprintf("%s.md", post.GenerateFileNamePrefix(postTitle))
	htmlFileName := fmt.Sprintf("%s.html", post.GenerateFileNamePrefix(postTitle))
	fmt.Println("going to publish: " + mdFileName)

	file, err := os.Open(constants.MarkdownPath + mdFileName)

	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := file.Stat()

	if err != nil {
		log.Fatal(err)
	}

	p := post.NewPostFromFile(fileInfo)

	p.Draft = false
	p.Date = time.Now()
	generatePost(p, htmlFileName)
	p.Update()

	generateIndex()
}

// generatePost is the helper function to actually create the .html file from the
// Post
func generatePost(post *post.Post, htmlFileName string) {
	fileName := constants.TemplatePath + "post.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(constants.HTMLPath + htmlFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	t.Execute(w, post)
	w.Flush()
}

// generateIndex loops through all of the Posts that have been published.
// Because the posts are already returned in descending date order, all we
// have to do is create the HTML
func generateIndex() {
	posts := stats.GetPublishedPosts()

	fileName := constants.TemplatePath + "index.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(constants.RootPath + "index.html")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	t.ExecuteTemplate(w, "header", nil)

	for _, post := range posts {
		t.ExecuteTemplate(w, "body", post)
	}

	t.ExecuteTemplate(w, "footer", nil)
	w.Flush()

}
