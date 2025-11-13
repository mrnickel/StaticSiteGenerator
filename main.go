package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Version is set at build time via ldflags
var version = "dev"

func main() {

	if len(os.Args) <= 1 || os.Args[1] == "help" {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "publish":
		// publish(os.Args[2])
		tmpP := NewPost(os.Args[2])

		fmt.Println("Publish the markdown file specified")
		file, err := os.Open(tmpP.MarkdownPath())
		if err != nil {
			log.Fatal(err)
		}
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		p, err := NewPostFromFile(fileInfo)
		if err != nil {
			log.Fatal(err)
		}
		err = p.Publish(false)
		if err != nil {
			log.Fatal(err)
		}
		return
	case "create":
		fmt.Println("Create a Post")
		p := NewPost(os.Args[2])
		p.Update()
		return
	case "stats":
		fmt.Println("Get the stats for this site")
		GetStats()
		return
	case "listdrafts":
		ListDrafts()
		return
	case "newsite":
		fmt.Println("TODO!")
		return
	case "preview":
		tmpP := NewPost(os.Args[2])
		file, err := os.Open(tmpP.MarkdownPath())
		if err != nil {
			log.Fatal(err)
		}
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		p, err := NewPostFromFile(fileInfo)
		if err != nil {
			log.Fatal(err)
		}
		err = p.Preview()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Launching your browser and go to http://localhost:8080/%s", p.HTMLPath())

		http.HandleFunc("/", makeGzipHandler(staticHandler))
		http.ListenAndServe(":8080", nil)
		openbrowser(fmt.Sprintf("http://localhost:8080/%s", p.HTMLPath()))
	case "standup":
		port := "8080" 
		
		if len(os.Args) > 2 {
			providedPort := os.Args[2]
			if portNum, err := strconv.Atoi(providedPort); err == nil && portNum > 0 && portNum < 65536 {
				port = providedPort
			} else {
				fmt.Printf("Invalid port '%s'. Using default port 8080.\n", providedPort)
			}
		}
		
		fmt.Printf("Now listening on port %s. Visit http://localhost:%s\n", port, port)
		http.HandleFunc("/", makeGzipHandler(staticHandler))
		http.ListenAndServe(":"+port, nil)
		openbrowser("http://localhost:" + port)
	case "regenerate":
		fmt.Println("Regenerating all published posts")

		posts := GetPublishedPosts()
		//sort posts by date descending
		for i := 0; i < len(posts)/2; i++ {
			posts[i], posts[len(posts)-1-i] = posts[len(posts)-1-i], posts[i]
		}

		for _, post := range posts {
			fmt.Printf("%s -- %s\n", post.Title(), post.Date())
			post.Publish(true)
		}

		return
	case "uuidify":
		fmt.Println("Adding GUIDs to all markdown posts that don't have one")
		uuidify()
		return
	case "version":
		fmt.Printf("Static Site Generator version %s\n", version)
		return
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("You must choose one of the following options:")
	fmt.Println("publish \"Blog Title here\"")
	fmt.Println("create \"Blog Title here\"")
	fmt.Println("stats (this will list stats about your site)")
	fmt.Println("listdrafts (this will list the titles of all your posts still in draft mode)")
	fmt.Println("preview \"Blog Title here\"")
	fmt.Println("newsite (creates a new site -- still needs to be implemented)")
	fmt.Println("standup [port] (this will start a web server that can handle requests, default port: 8080)")
	fmt.Println("regenerate (this will regenerate all published pages with the new templates)")
	fmt.Println("uuidify (adds GUIDs to all markdown posts that don't have one)")
	fmt.Println("version (display version information)")
}

// func publish(file string) {
// 	tmpP := NewPost(file)

// 	fmt.Println("Publish the markdown file specified")
// 	file, err := os.Open(tmpP.MarkdownPath())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	p, err := NewPostFromFile(fileInfo)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = p.Publish()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {

	if len(r.URL.Path) <= 1 {
		r.URL.Path = "/index.html"
		// fmt.Println("heh")
		// return
	}

	file, err := os.Open(r.URL.Path[1:])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileStat, err := os.Stat(r.URL.Path[1:])
	if err != nil {
		panic(err)
	}

	_, filename := path.Split(r.URL.Path[1:])
	t := fileStat.ModTime()
	http.ServeContent(w, r, filename, t, file)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

// uuidify loops through all markdown files and adds a GUID to posts that don't have one
func uuidify() {
	files, err := os.ReadDir(MarkdownPath)
	if err != nil {
		log.Fatal(err)
	}

	updated := 0
	skipped := 0

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			log.Printf("Error reading file info for %s: %v\n", file.Name(), err)
			continue
		}

		post, err := NewPostFromFile(fileInfo)
		if err != nil {
			log.Printf("Error parsing post %s: %v\n", file.Name(), err)
			continue
		}

		// Check if post already has a GUID
		if post.GUID() != "" {
			fmt.Printf("✓ Skipping %s (already has GUID: %s)\n", post.Title(), post.GUID())
			skipped++
			continue
		}

		// Generate and set a new GUID
		newGUID := GenerateGUID()
		post.SetGUID(newGUID)
		
		fmt.Printf("+ Adding GUID to %s: %s\n", post.Title(), newGUID)
		post.Update()
		updated++
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Updated: %d posts\n", updated)
	fmt.Printf("  Skipped: %d posts\n", skipped)
}
