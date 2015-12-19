# Static Site Generator

This project essentially has 2 things that it does

1. It creates a markdown file that we can add to our github repo. This markdown file will include the following in the header:

	1. Title
	2. Published Date
	3. State (draft/published)

2. It parses through all of the markdown files creating static html versions of it, adds them to the code repo, and commits them


# You probably don't want to use this

This was created by myself as a fun project in order to get my feet wet with Go. If you want a real static site generator I suggest you use [Hugo](https://gohugo.io/).