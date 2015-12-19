package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/knieriem/markdown"
)

// Post is a post that we can do stuff with
type Post struct {
	Date        time.Time
	Draft       bool
	Title       string
	MDContent   string
	HTMLContent string
	Summary     string
}

// PostsByDate is a descending sortable slice of Post structs
type PostsByDate []*Post

func (slice PostsByDate) Len() int {
	return len(slice)
}

func (slice PostsByDate) Less(i, j int) bool {
	return slice[i].Date.After(slice[j].Date)
}

func (slice PostsByDate) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// NewPost creates a new Post struct and defaults the time
// to right now, the draft value to true and the title to the
// specified title
func NewPost(title string) *Post {
	p := new(Post)
	p.Date = time.Now()
	p.Draft = true
	p.Title = title

	return p
}

// NewPostFromFile will create a new post based on the file
// specified
func NewPostFromFile(fileInfo os.FileInfo) *Post {
	// p := new(Post)

	var p *Post

	file, err := os.Open(baseArticlePath + fileInfo.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if scanner.Text() == "+++" {
			p, err = parseHeader(scanner)
			if err != nil {
				log.Fatal(err)
			}
		}

		if p.MDContent == "" {
			p.MDContent = scanner.Text()
		} else {
			p.MDContent = p.MDContent + "\n" + scanner.Text()
		}
	}

	p.HTMLContent = convertContentToHTML(p)

	return p
}

// Update will update the .md file associated with this post. Typically
// done after the post is published
func (p *Post) Update() {
	file, err := os.Create(baseArticlePath + fmt.Sprintf("%s.md", strings.Replace(strings.ToLower(p.Title), " ", "+", -1)))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(p.GetString())
}

// parseHeader is a helper function for NewPostFromFile. It will
// parse out the text between the start +++ and end +++ tags.
//
// Returns the beginning of the Post struct
func parseHeader(scanner *bufio.Scanner) (*Post, error) {
	p := new(Post)
	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "date") {
			p.Date, _ = parseDate(scanner.Text())
		} else if strings.Contains(scanner.Text(), "draft") {
			p.Draft = parseDraft(scanner.Text())
		} else if strings.Contains(scanner.Text(), "title") {
			p.Title = parseTitle(scanner.Text())
		}

		if scanner.Text() == "+++" {
			scanner.Scan() // get us past the +++ line
			return p, nil
		}
	}

	return nil, errors.New("Cant parse file")
}

// parseDate will parse the date string into a proper time.Time
// variable. If we fail at parsing it, we return an error
func parseDate(dateLine string) (time.Time, error) {
	dateStr := strings.TrimPrefix(dateLine, "date = ")
	return time.Parse(time.RFC3339, dateStr)
}

// parseDraft will tell you if this Post is draft = true, or
// draft = false
func parseDraft(draftLine string) bool {
	if strings.TrimPrefix(draftLine, "draft = ") == "true" {
		return true
	}
	return false
}

// parseTitle simply returns the title value.
//
// Used the helper function here in order to keep things consistent
func parseTitle(titleLine string) string {
	return strings.TrimPrefix(titleLine, "title = ")
}

// convertPostToHtml will convert the Post.content markdown
// into an HTML file
func convertContentToHTML(p *Post) string {
	reader := bytes.NewReader([]byte(p.MDContent))
	parser := markdown.NewParser(&markdown.Extensions{Smart: true})
	dst := new(bytes.Buffer)
	parser.Markdown(reader, markdown.ToHTML(dst))

	return dst.String()
}

// GetString returns the posts string value that we would
// potentially write out to a file
func (p *Post) GetString() string {
	var buffer bytes.Buffer
	buffer.WriteString("+++\n")
	buffer.WriteString(fmt.Sprintf("date = %s\n", p.Date.Format(time.RFC3339)))

	if p.Draft {
		buffer.WriteString("draft = true\n")
	} else {
		buffer.WriteString("draft = false\n")
	}

	buffer.WriteString("title = " + p.Title + "\n")
	buffer.WriteString("+++\n")
	buffer.WriteString("\n")

	if p.MDContent == "" {
		buffer.WriteString("# " + p.Title + "\n")
	} else {
		buffer.WriteString(p.MDContent)
	}

	return buffer.String()
}
