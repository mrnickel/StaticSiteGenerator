package post

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mrnickel/StaticSiteGenerator/config"
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

// Update will update the .md file associated with this post. Typically
// done after the post is published
func (p *Post) Update() {
	file, err := os.Create(config.MarkdownPath + fmt.Sprintf("%s.md", strings.Replace(strings.ToLower(p.Title), " ", "+", -1)))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(p.String())
}

// GetString returns the posts string value that we would
// potentially write out to a file
func (p *Post) String() string {
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
