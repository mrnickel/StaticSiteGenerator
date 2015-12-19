# Static Site Generator

This project essentially has 2 things that it does

1. It creates a markdown file that we can add to our github repo. This markdown file will include the following in the header:

	1. Date
	2. Draft
	3. Title

2. It parses through all of the markdown files creating static html versions of it

# You probably don't want to use this

This was created by myself as a fun project in order to get my feet wet with Go. If you want a real static site generator I suggest you use [Hugo](https://gohugo.io/).

# TODO
- [ ] Update this readme to be more informative
- [ ] Write basic wiki in order to remind myself how to use it, for the 2-3 blog posts I actually make / year
- [ ] properly figure out packages, so my helper / private functions don't pollute the project
- [ ] become more consistent with naming... i.e. article vs post