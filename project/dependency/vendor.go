package dependency

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

func (d Dependency) Vendor(vendorRoot string, force bool) error {
	path := filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.rootPackage))
	if force {
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.MkdirAll(path, 0775); err != nil {
		return err
	}
	v := modules.GetVCS(d.VCS)
	if v == nil {
		return fmt.Errorf("Unkown Version Control System : %s", d.VCS)
	}
	_, err := vcs.CloneAtRevision(v, d.Remote, d.Revision, path)
	return err
}
