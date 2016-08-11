package project

import "strings"

//TODO Function to get project's root directory

type empty struct{}

//Project represents a Go project
type Project struct {
	Root string
}

var (
	pseudoPackages = []string{
		"C",
	}
)

func in(elem string, arr []string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

//New creates a new `Project`
func New(root string) *Project {
	return &Project{
		Root: root,
	}
}

//TODO Only walk through root directory and follow dependencies
//TODO Clean this function

//GetDependencies gets all dependencies of the project
func (p Project) GetDependencies() (<-chan *Dependency, <-chan error) {
	out := make(chan *Dependency)
	errch := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errch)
		directories := []string{p.Root}
		m := make(map[string]empty) //Set of dependencies
		for it := 0; it < len(directories); it++ {
			var (
				filepaths <-chan string
			)
			if it == 0 { //Walk through project directory, without filtering test files
				filepaths = walk(directories[it], isValidFile, isValidDir, sendError(errch))
			} else { //Walk though other directories, skipping subdirectories
				filepaths = walk(directories[it], isValidGoFile, skipDirExcept(directories[it]), sendError(errch))
			}
			for file := range filepaths {
				imports, err := GetImports(file)
				if err != nil {
					errch <- NewErrorMessage(err).WithField("file", file).WithMessage("Could not get imports from file")
					return
				}
				for _, i := range imports {
					if !in(GetName(i), pseudoPackages) { //Check if this import is not a pseudo package
						s, err := NewDependency(i)
						if err != nil {
							errch <- NewErrorMessage(err).WithField("import", i.Path.Value).
								WithField("directory", directories[it]).
								WithMessage("Could create new import")
							return
						} else if _, ok := m[i.Path.Value]; !s.Package.Goroot && !ok { //Filter packages from stdlib
							m[i.Path.Value] = empty{}
							directories = append(directories, s.Package.Dir)
							if !strings.HasPrefix(s.Package.Dir, p.Root) { //Do not send imports from the same project
								out <- s
							}
						}
					}
				}
			}
		}
	}()
	return out, errch
}
