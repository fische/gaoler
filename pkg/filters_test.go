package pkg_test

import (
	"os"
	"time"

	. "github.com/fische/gaoler/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	return f.mode
}

func (f fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f fileInfo) IsDir() bool {
	return f.isDir
}

func (f fileInfo) Sys() interface{} {
	return f.sys
}

var _ = Describe("Filters", func() {
	var (
		file  fileInfo
		name  string
		isDir bool
	)

	JustBeforeEach(func() {
		file = fileInfo{
			name:  name,
			isDir: isDir,
		}
	})

	Describe("IsNotGoTestFile", func() {
		Context("When file is not a Go test", func() {
			BeforeEach(func() {
				name = "notATest.go"
			})

			It("should return true", func() {
				Expect(IsNotGoTestFile(file)).To(BeTrue())
			})
		})

		Context("When file has the Go test file extension", func() {
			BeforeEach(func() {
				name = "is_a_test.go"
			})
			Context("When file is a directory", func() {
				BeforeEach(func() {
					isDir = true
				})

				It("should return true", func() {
					Expect(IsNotGoTestFile(file)).To(BeTrue())
				})
			})

			Context("When file is not a directory", func() {
				BeforeEach(func() {
					isDir = false
				})

				It("should return false", func() {
					Expect(IsNotGoTestFile(file)).To(BeFalse())
				})
			})
		})

		Context("When file is called testdata", func() {
			BeforeEach(func() {
				name = "testdata"
			})
			Context("When file is a directory", func() {
				BeforeEach(func() {
					isDir = true
				})

				It("should return false", func() {
					Expect(IsNotGoTestFile(file)).To(BeFalse())
				})
			})

			Context("When file is not a directory", func() {
				BeforeEach(func() {
					isDir = false
				})

				It("should return true", func() {
					Expect(IsNotGoTestFile(file)).To(BeTrue())
				})
			})
		})
	})
})
