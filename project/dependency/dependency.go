package dependency

import (
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type Dependency struct {
	rootPackage string

	repository vcs.Repository
	VCS        string `json:",omitempty" yaml:",omitempty"`
	Revision   string `json:",omitempty" yaml:",omitempty"`
	Remote     string `json:",omitempty" yaml:",omitempty"`
	Branch     string `json:",omitempty" yaml:",omitempty"`

	Packages []*pkg.Package
}

func New(p *pkg.Package) *Dependency {
	return &Dependency{
		rootPackage: p.Path(),
		Packages:    []*pkg.Package{p},
	}
}

func (d *Dependency) Add(p *pkg.Package) (added bool) {
	if d.HasPackage(p.Path()) {
		return
	}
	added = true
	d.Packages = append(d.Packages, p)
	return
}

func (d *Dependency) SetRootPackage(rootPackage string) {
	d.rootPackage = rootPackage
}

func (d Dependency) HasPackage(packagePath string) bool {
	for _, pkg := range d.Packages {
		if pkg.Path() == packagePath {
			return true
		}
	}
	return false
}

func (d Dependency) IsVendored() bool {
	for _, p := range d.Packages {
		if !p.IsVendored() {
			return false
		}
	}
	return true
}

func (d *Dependency) OpenRepository(p *pkg.Package) error {
	var (
		err  error
		path string

		repo        vcs.Repository
		rootPackage string
		remote      string
		revision    string
	)
	if repo, err = modules.OpenRepository(p.Dir()); err != nil {
		return err
	} else if path, err = repo.GetPath(); err != nil {
		return err
	} else if rootPackage, err = pkg.GetPackagePath(path); err != nil {
		return err
	} else if remote, err = repo.GetRemote(); err != nil {
		return err
	} else if revision, err = repo.GetRevision(); err != nil {
		return err
	}
	d.VCS = repo.GetVCSName()
	d.repository = repo
	d.rootPackage = rootPackage
	d.Remote = remote
	d.Revision = revision
	return nil
}

func (d *Dependency) Import(srcPath string, flags pkg.Flags) error {
	for _, p := range d.Packages {
		if err := p.Import(srcPath, flags); err != nil {
			return err
		}
	}
	return nil
}

func (d Dependency) RootPackage() string {
	return d.rootPackage
}

func (d Dependency) IsOpened() bool {
	return d.repository != nil
}

func (d Dependency) IsVendorable() bool {
	return d.VCS != "" && d.Remote != "" && d.Revision != ""
}
