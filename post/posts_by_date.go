package post

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
