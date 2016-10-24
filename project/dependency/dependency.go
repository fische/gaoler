package dependency

import (
	"errors"
	"strings"

	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type Dependency struct {
	rootPackage string
	repository  vcs.Repository
	*pkg.Set
	*State `json:",omitempty" yaml:",omitempty"`
}

func New(p *pkg.Package) *Dependency {
	s := pkg.NewSet()
	s.Insert(p, true)
	return &Dependency{
		rootPackage: p.Path(),
		Set:         s,
	}
}

func (d *Dependency) Add(p *pkg.Package) (added bool) {
	added = d.Set.Insert(p, false)
	if added {
		if strings.HasPrefix(d.rootPackage, p.Path()) && d.rootPackage != p.Path() {
			d.rootPackage = p.Path()
		}
	}
	return
}

func (d *Dependency) SetRootPackage() error {
	if d.repository == nil {
		return errors.New("Package repository has not been opened.")
	}
	var (
		path        string
		rootPackage string
		err         error
	)
	if path, err = d.repository.GetPath(); err != nil {
		return err
	} else if rootPackage, err = pkg.GetPackagePath(path); err != nil {
		return err
	}
	d.rootPackage = rootPackage
	return nil
}

func (d *Dependency) OpenRepository(dir string) error {
	var (
		repo vcs.Repository
		err  error
	)
	if repo, err = modules.OpenRepository(dir); err != nil {
		return err
	}
	d.repository = repo
	return nil
}

func (d Dependency) RootPackage() string {
	return d.rootPackage
}

func (d Dependency) Repository() vcs.Repository {
	return d.repository
}

func (d Dependency) IsVendored() bool {
	for _, p := range d.Set.Packages() {
		if !p.IsVendored() {
			return false
		}
	}
	return true
}
