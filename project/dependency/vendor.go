package dependency

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

func (d *Dependency) Vendor(vendorRoot string) error {
	if d.State == nil {
		return errors.New("State has not been locked.")
	}
	var (
		err  error
		path = filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.rootPackage))
	)

	if err = os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	} else if err = os.MkdirAll(path, 0775); err != nil {
		return err
	}
	v := modules.GetVCS(d.vcs)
	if v == nil {
		return fmt.Errorf("Unkown Version Control System : %s", d.vcs)
	}
	d.repository, err = vcs.CloneAtRevision(v, d.remote, d.revision, path)
	return err
}

func (d *Dependency) Update(vendorRoot string) (updated bool, err error) {
	var revision string
	if err = d.Vendor(vendorRoot); err != nil {
		return
	} else if err = d.repository.CheckoutBranch(d.branch); err != nil {
		return
	} else if revision, err = d.repository.GetRevision(); err != nil {
		return
	}
	if d.revision != revision {
		updated = true
		d.revision = revision
	}
	return
}
