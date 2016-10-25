package dependency

import (
	"fmt"
	"strings"

	"github.com/fische/gaoler/pkg"
)

type DiffDependency struct {
	added     *pkg.Set
	removed   *pkg.Set
	untouched *pkg.Set
}

type DiffSet map[string]*DiffDependency

func (d *Dependency) Apply(diff *DiffDependency, onPackageAdded func(p *pkg.Package, dep *Dependency) error) error {
	for packagePath, p := range diff.added.Packages() {
		d.Insert(p, true)
		if strings.HasPrefix(d.rootPackage, packagePath) &&
			d.rootPackage != packagePath {
			d.rootPackage = packagePath
		}
		if onPackageAdded != nil {
			if err := onPackageAdded(p, d); err != nil {
				return err
			}
		}
	}
	for packagePath := range diff.removed.Packages() {
		d.Remove(packagePath)
	}
	return nil
}

func (d Dependency) Diff(o *pkg.Set) *DiffDependency {
	diff := &DiffDependency{
		added:     pkg.NewSet(),
		removed:   pkg.NewSet(),
		untouched: pkg.NewSet(),
	}
	for packagePath, p := range d.Packages() {
		if o.Has(packagePath) {
			diff.untouched.Insert(p, false)
			o.Remove(packagePath)
		} else {
			diff.removed.Insert(p, false)
		}
	}
	for packagePath, p := range o.Packages() {
		if strings.HasPrefix(packagePath, d.rootPackage) ||
			strings.HasPrefix(d.rootPackage, packagePath) {
			diff.added.Insert(p, false)
			o.Remove(packagePath)
		}
	}
	return diff
}

func (s *Set) Apply(diff DiffSet) error {
	for rootPackage, diffDep := range diff {
		if diffDep.IsRemoved() {
			delete(s.dependencies, rootPackage)
			continue
		}
		dep, ok := s.dependencies[rootPackage]
		if !ok {
			return fmt.Errorf("Could not retrieve dependency %s", rootPackage)
		} else if err := dep.Apply(diffDep, s.OnPackageAdded); err != nil {
			return err
		} else if rootPackage != dep.rootPackage {
			delete(s.dependencies, rootPackage)
			s.dependencies[dep.rootPackage] = dep
		}
	}
	return nil
}

func (s Set) Diff(o *pkg.Set) DiffSet {
	diffs := make(DiffSet)
	for rootPackage, dep := range s.dependencies {
		diffs[rootPackage] = dep.Diff(o)
	}
	return diffs
}

func (d DiffDependency) IsRemoved() bool {
	return len(d.added.Packages()) == 0 && len(d.untouched.Packages()) == 0
}

func (d DiffDependency) Added() *pkg.Set {
	return d.added
}

func (d DiffDependency) Removed() *pkg.Set {
	return d.removed
}

func (d DiffDependency) Untouched() *pkg.Set {
	return d.untouched
}
