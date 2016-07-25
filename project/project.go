package project

//TODO Function to get project's root directory

//Project represents a Go project
type Project struct {
	Root string
}

//New creates a new `Project`
func New(root string) *Project {
	return &Project{
		Root: root,
	}
}
