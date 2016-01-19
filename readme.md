# You probably don't want to use this
This was created by myself as a fun project in order to get my feet wet with Go. If you want a real static site generator I suggest you use [Hugo](https://gohugo.io/).

# Static Site Generator
This project essentially has 2 things that it does

1. It creates a markdown file that we can add to our github repo. This markdown file will include the following in the header:

	1. Date
	2. Draft
	3. Title

2. It parses through all of the markdown files creating static html versions of it

# Commands
- `StaticSiteGenerator publish "*title of post*"`
- `StaticSiteGenerator create "*title of post*"`
- `StaticSiteGenerator stats`
- `StaticSiteGenerator listdrafts`

#Requirements
In order to use this VERY basic system you will have to follow a few basic rules:

## Folder structure
We assume a very specific folder structure:
- / _Root path to your site_
	- /html _This is where the generated static HTML files will go_
	- /md _This is where `StaticSiteGenerator create` will put the generated markdown files_
	- /templates _This is where the template files go.

##Templates
The templates are pretty basic and utilize Go's built in [template package](https://golang.org/pkg/html/template/). All of the placeholders are the Post struct's fields

```
type Post struct {
	Date           time.Time
	Draft          bool
	Title          string
	MDContent      string
	HTMLContent    string
	Summary        string
	FileNamePrefix string
}
```

# TODO
- [x] Update this readme to be more informative
- [x] Write basic wiki in order to remind myself how to use it, for the 2-3 blog posts I actually make / year
- [x] properly figure out packages, so my helper / private functions don't pollute the project
- [x] become more consistent with naming... i.e. article vs post
- [x] make constants generic so as to presume the app is installed and not being run via `go run`
- [x] document website requirements. i.e. folder structure and template structure
- [x] post to reddit / hackernews / twitter so as to get feedback
- [ ] create *preview* function to see what my markdown looks like when in the templates
- [ ] vendoring?
- [ ] RSS feed generation
- [x] refactor code to be based around post struct
- [x] put everything under root as outlined in Peter Bourgon's [Go: Best Practices for Production Environments](http://peter.bourgon.org/go-in-production/)
- [ ] create **newsite** command
- [ ] better check user data from start so as to not panic and quit right off the bat (better error handling)
