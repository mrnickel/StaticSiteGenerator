package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

// Publish will take the Post.Title of the item we wish to publish, generate the HTML from the
// template, update the .md file's draft flag to false and re-generate the index page.
// Maybe one day I'll add other flags, such as "+tweet" in order to connect to twitter
// and post on my behalf
func Publish(post string) {
	mdFileName := fmt.Sprintf("%s.md", strings.Replace(strings.ToLower(post), " ", "+", -1))
	htmlFileName := fmt.Sprintf("%s.html", strings.Replace(strings.ToLower(post), " ", "+", -1))
	fmt.Println("going to publish: " + mdFileName)

	file, err := os.Open(baseMarkdownPath + mdFileName)

	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := file.Stat()

	if err != nil {
		log.Fatal(err)
	}

	p := NewPostFromFile(fileInfo)

	p.Draft = false
	p.Date = time.Now()
	generatePost(p, htmlFileName)
	p.Update()

	generateIndex()
}

// generatePost is the helper function to actually create the .html file from the
// Post
func generatePost(post *Post, htmlFileName string) {
	fileName := baseTemplatePath + "post.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(baseHTMLPath + htmlFileName)
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
	posts := getPublishedPosts()

	fileName := baseTemplatePath + "index.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(baseRootPath + "index.html")
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
