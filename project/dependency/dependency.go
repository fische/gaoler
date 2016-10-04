package dependency

import (
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/vcs"
)

type Dependency struct {
	RootPackage string `json:"-" yaml:"-"`

	Repository vcs.Repository `json:"-" yaml:"-"`
	VCS        string         `json:",omitempty" yaml:",omitempty"`
	Revision   string         `json:",omitempty" yaml:",omitempty"`
	Remote     string         `json:",omitempty" yaml:",omitempty"`
	Branch     string         `json:",omitempty" yaml:",omitempty"`

	Packages []*pkg.Package
}

func New(p *pkg.Package) *Dependency {
	return &Dependency{
		RootPackage: p.Path,
		Packages:    []*pkg.Package{p},
	}
}

func (d *Dependency) Add(p *pkg.Package) (added bool) {
	if d.HasPackage(p.Path) {
		return
	}
	added = true
	d.Packages = append(d.Packages, p)
	return
}

func (d Dependency) IsVendored() bool {
	for _, p := range d.Packages {
		if !p.Vendored {
			return false
		}
	}
	return true
}

func (d Dependency) HasPackage(packagePath string) bool {
	for _, pkg := range d.Packages {
		if pkg.Path == packagePath {
			return true
		}
	}
	return false
}

func (d Dependency) IsVendorable() bool {
	return d.VCS != "" && d.Remote != "" && d.Revision != ""
}

func (d Dependency) IsUpdatable() bool {
	return d.IsVendorable() && d.Branch != ""
}
