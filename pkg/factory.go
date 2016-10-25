package pkg

type Factory interface {
	New(packagePath string) (p *Package, nextDirectories []string, err error)
}

type BasicFactory struct {
	SrcPath              string
	IgnoreVendor         bool
	ImportCanFail        bool
	IncludePseudoPackage bool
}

func (f BasicFactory) New(packagePath string) (p *Package, nextDirectories []string, err error) {
	if IsPseudoPackage(packagePath) {
		if f.IncludePseudoPackage {
			p = New(packagePath)
		}
		return
	}
	p, err = Import(packagePath, f.SrcPath, f.IgnoreVendor)
	if err != nil {
		if f.ImportCanFail {
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
