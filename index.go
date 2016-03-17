package main

// Index is a struct that allows us to send in all data
// required for rendering each paginated index page
type Index struct {
	Posts        []Post
	HasNext      bool
	NextPage     string
	HasPrevious  bool
	PreviousPage string
}
