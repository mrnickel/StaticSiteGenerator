package post

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mrnickel/StaticSiteGenerator/constants"
)

// Post is a post that we can do stuff with
type Post struct {
	Date           time.Time
	Draft          bool
	Title          string
	MDContent      string
	HTMLContent    string
	Summary        string
	FileNamePrefix string
}

// NewPost creates a new Post struct and defaults the time
// to right now, the draft value to true and the title to the
// specified title
func NewPost(title string) *Post {
	p := new(Post)
	p.Date = time.Now()
	p.Draft = true
	p.Title = title
	p.FileNamePrefix = GenerateFileNamePrefix(title)

	return p
}

// GenerateFileNamePrefix will return the prefix of the file name
// i.e. if the Title is "This is my post", it will return "this+is+my+post"
// so that we can append .html or .md to the filename and open the appropriate
// file
func GenerateFileNamePrefix(title string) string {
	return strings.Replace(strings.ToLower(title), " ", "_", -1)
}

// Update will update the .md file associated with this post. Typically
// done after the post is published
func (p *Post) Update() {
	file, err := os.Create(constants.MarkdownPath + fmt.Sprintf("%s.md", strings.Replace(strings.ToLower(p.Title), " ", "_", -1)))
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

	buffer.WriteString(p.MDContent)

	return buffer.String()
}
