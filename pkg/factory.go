package pkg

type Factory interface {
	New(packagePath string) (p *Package, nextDirectories []string, err error)
}

type DefaultFactory struct {
	SrcPath             string
	IgnoreVendor        bool
	ImportCanFail       bool
	ReturnPseudoPackage bool
}

func (d DefaultFactory) New(packagePath string) (p *Package, nextDirectories []string, err error) {
	if IsPseudoPackage(packagePath) {
		if d.ReturnPseudoPackage {
			p = New(packagePath)
		}
		return
	}
	p, err = Import(packagePath, d.SrcPath, d.IgnoreVendor)
	if err != nil {
		if d.ImportCanFail {
			err = nil
		}
		return
	} else if p.IsStandardPackage() {
		p = nil
		return
	}
	if p.Dir() != "" {
		nextDirectories = []string{p.Dir()}
	}
	if p.IsLocal() {
		p = nil
		return
	}
	return
}
