package types

// FBEntity ...
type FBEntity struct {
	isFile bool

	path         string
	title        string
	extension    string
	modifiedDate string
	size         string
}
