package post

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/knieriem/markdown"
	"github.com/mrnickel/StaticSiteGenerator/config"
)

// NewPostFromFile will create a new post based on the file
// specified
func NewPostFromFile(fileInfo os.FileInfo) *Post {

	var p *Post

	file, err := os.Open(config.MarkdownPath + fileInfo.Name())
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
