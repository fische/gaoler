package project

//TODO Function to get project's root directory

type empty struct{}

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

//GetDependencies gets all dependencies of the project
func (p Project) GetDependencies() (<-chan *Import, <-chan error) {
	out := make(chan *Import)
	errch := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errch)
		directories := []string{p.Root}
		m := make(map[string]empty)
		for _, directory := range directories {
			filepaths := walk(directory, isValidFile, isValidDir, sendError(errch))
			for file := range filepaths {
				imports, err := GetImports(file)
				if err != nil {
					errch <- NewErrorMessage(err).WithField("file", file).WithMessage("Could not get imports from file")
					return
				}
				for _, i := range imports {
					s, err := NewImport(i)
					if err != nil {
						errch <- NewErrorMessage(err).WithField("import", i.Path.Value).WithMessage("Could create new import")
						return
					} else if _, ok := m[i.Path.Value]; !s.Goroot && !ok {
						m[i.Path.Value] = empty{}
						directories = append(directories, s.Dir)
						out <- s
					}
				}
			}
		}
	}()
	return out, errch
}
