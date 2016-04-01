package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/russross/blackfriday"
)

const (
	// MarkdownPath is the direcotry path to the written markdown files
	MarkdownPath = "md"
	// TemplatePath is the directory path to the html templates
	TemplatePath = "templates"
	// HTMLPath is the directory path to where the generated HTML files will go
	HTMLPath = "html"
	// RootPath is the directory path to the root of the website
	RootPath = ""
	// PageSize is the number of posts in a page
	PageSize = 10
)

// Post is an interface ontop of the post struct
type Post interface {
	Date() time.Time
	Draft() bool
	Title() string
	MarkdownContent() string
	HTMLContent() string
	MarkdownPath() string
	HTMLPath() string
	Publish() error
	Update()
	String() string
	Preview() error
	Summary() string
}

// An array of posts that we can use in the templates
// type Context struct {
// 	Posts []Post
// }

// post is a post that we can do stuff with
type post struct {
	date            time.Time
	draft           bool
	title           string
	markdownContent string
}

// PostsByDate is a descending sortable slice of Post structs
type PostsByDate []Post

func (slice PostsByDate) Len() int {
	return len(slice)
}

func (slice PostsByDate) Less(i, j int) bool {
	return slice[i].Date().After(slice[j].Date())
}

func (slice PostsByDate) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Date returns the post structs private date field
func (p *post) Date() time.Time {
	return p.date
}

// Draft returns the post structs private draft field
func (p *post) Draft() bool {
	return p.draft
}

// Title returns the post structs private title field
func (p *post) Title() string {
	return p.title
}

// MarkdownContent returns the post structs private markdownContent field
func (p *post) MarkdownContent() string {
	return p.markdownContent
}

// HTMLContent returns the post structs private hTMLContent field
func (p *post) HTMLContent() string {
	html := blackfriday.MarkdownCommon([]byte(p.markdownContent))
	return string(html[:])
}

func (p *post) filePrefix() string {
	return strings.Replace(strings.ToLower(p.title), " ", "_", -1)
}

func (p *post) MarkdownPath() string {
	return fmt.Sprintf("%s/%s.md", MarkdownPath, p.filePrefix())
}

func (p *post) SetMarkdownContent(content string) {
	p.markdownContent = content
}

func (p *post) HTMLPath() string {
	return fmt.Sprintf("%s/%s.html", HTMLPath, p.filePrefix())
}

// String returns the posts string value that we would
// potentially write out to a file
func (p *post) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("---\n")
	buffer.WriteString(fmt.Sprintf("date: %s\n", p.date.Format(time.RFC3339)))

	if p.draft {
		buffer.WriteString("draft: true\n")
	} else {
		buffer.WriteString("draft: false\n")
	}

	buffer.WriteString("title: " + p.title + "\n")
	buffer.WriteString("---\n")
	buffer.WriteString("\n")

	buffer.WriteString(p.markdownContent)

	return buffer.String()
}

// Returns the first paragraph in the post
func (p *post) Summary() string {
	paragraphs := strings.Split(p.MarkdownContent(), "\n")

	if len(paragraphs) == 0 {
		panic("shit")
	}

	html := blackfriday.MarkdownCommon([]byte(paragraphs[0]))
	return string(html[:])
}

// Update will update the .md file associated with this post. Typically
// done after the post is published
func (p *post) Update() {
	file, err := os.Create(p.MarkdownPath())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(p.String())
}

// NewPost creates a new Post struct and defaults the time
// to right now, the draft value to true and the title to the
// specified title
func NewPost(title string) Post {
	p := new(post)
	p.date = time.Now()
	p.draft = true
	p.title = title

	return p
}

// NewPostFromFile will create a new post based on the file
// if the file doesn't exist then we return an error
func NewPostFromFile(fileInfo os.FileInfo) (Post, error) {
	tmpPost := new(post)
	tmpPost.draft = false

	file, err := os.Open(MarkdownPath + "/" + fileInfo.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	parsingHeader := false

	for scanner.Scan() {
		if scanner.Text() == "---" && !parsingHeader {
			parsingHeader = true
		} else if parsingHeader && strings.HasPrefix(scanner.Text(), "date") {
			dateStr := strings.TrimPrefix(scanner.Text(), "date: ")
			tmpPost.date, _ = time.Parse(time.RFC3339, dateStr)
		} else if parsingHeader && strings.HasPrefix(scanner.Text(), "draft") {
			if strings.TrimPrefix(scanner.Text(), "draft: ") == "true" {
				tmpPost.draft = true
			}
		} else if parsingHeader && strings.HasPrefix(scanner.Text(), "title") {
			tmpPost.title = strings.TrimPrefix(scanner.Text(), "title: ")
		} else if parsingHeader && scanner.Text() == "---" {
		} else if !parsingHeader && tmpPost.MarkdownContent() == "" {
			tmpPost.markdownContent = scanner.Text()
		} else {
			tmpPost.markdownContent = tmpPost.MarkdownContent() + "\n" + scanner.Text()
		}
	}

	tmpPost.markdownContent = strings.TrimSpace(tmpPost.MarkdownContent())

	return tmpPost, nil
}

// Publish generate the HTML from the template, update the .md file's draft flag to false
// and re-generate the index page.
// Maybe one day I'll add other flags, such as "+tweet" in order to connect to twitter
// and post on my behalf
func (p *post) Publish() error {
	p.draft = false
	p.date = time.Now()
	p.Update()

	fileName := TemplatePath + "/post.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		return err
	}

	f, err := os.Create(p.HTMLPath())
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	t.Execute(w, p)
	w.Flush()

	err = generateIndex()
	if err != nil {
		return err
	}

	err = generateRss()
	if err != nil {
		return err
	}
	// now we gotta do the index page
	return nil
}

func (p *post) Preview() error {
	fileName := TemplatePath + "/post.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		return err
	}

	f, err := os.Create(p.HTMLPath())
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	t.Execute(w, p)
	w.Flush()

	return nil

}

// generateIndex loops through all of the Posts that have been published.
// Because the posts are already returned in descending date order, all we
// have to do is create the HTML
func generateIndex() error {
	posts := GetPublishedPosts()

	fileName := TemplatePath + "/index.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		return err
	}

	numItems := len(posts)

	numPages := int(math.Ceil(float64(numItems) / float64(PageSize)))

	for page := 1; page <= numPages; page++ {
		startIndex := (page - 1) * PageSize
		endIndex := startIndex + PageSize
		if endIndex > len(posts) {
			endIndex = len(posts)
		}

		indexPage := "index.html"

		if page > 1 {
			indexPage = fmt.Sprintf("index-%d.html", page)
		}

		f, err := os.Create(RootPath + indexPage)
		if err != nil {
			return err
		}

		defer f.Close()

		w := bufio.NewWriter(f)
		indexData := new(Index)
		indexData.Posts = posts[startIndex:endIndex]

		indexData.HasNext = false
		indexData.HasPrevious = false
		if page < numPages {
			indexData.HasNext = true
			indexData.NextPage = fmt.Sprintf("index-%d", page+1)
		}
		if page > 1 {
			indexData.HasPrevious = true
			indexData.PreviousPage = "index"

			if page > 2 {
				indexData.PreviousPage = fmt.Sprintf("index-%d", page-1)
			}

		}
		t.Execute(w, indexData)

		w.Flush()
	}

	return nil
}

func generateRss() error {
	posts := GetPublishedPosts()

	fileName := TemplatePath + "/rss.tmpl"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		return err
	}

	f, err := os.Create(RootPath + "rss.xml")
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	endIndex := PageSize
	if endIndex > len(posts) {
		endIndex = len(posts) - 1
	}

	t.Execute(w, posts[0:endIndex])

	w.Flush()
	return nil
}
