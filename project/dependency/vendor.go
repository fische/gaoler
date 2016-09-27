package dependency

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

func (d *Dependency) Vendor(vendorRoot string) error {
	var (
		err  error
		path = filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.rootPackage))
	)

	if err = os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	} else if err = os.MkdirAll(path, 0775); err != nil {
		return err
	}
	v := modules.GetVCS(d.VCS)
	if v == nil {
		return fmt.Errorf("Unkown Version Control System : %s", d.VCS)
	}
	d.repository, err = vcs.CloneAtRevision(v, d.Remote, d.Revision, path)
	return err
}

func (d *Dependency) Update(vendorRoot string) (update bool, err error) {
	var revision string
	if err = d.Vendor(vendorRoot); err != nil {
		return
	} else if err = d.repository.Checkout(d.Branch); err != nil {
		return
	} else if revision, err = d.repository.GetRevision(); err != nil {
		return
	}
	if d.Revision != revision {
		update = true
		d.Revision = revision
	}
	return
}
