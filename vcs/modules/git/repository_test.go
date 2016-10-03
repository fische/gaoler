package git_test

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fische/gaoler/vcs"
	. "github.com/fische/gaoler/vcs/modules/git"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	Describe("InitRepository", func() {
		var (
			r   vcs.Repository
			err error

			path string
			bare bool
		)

		JustBeforeEach(func() {
			r, err = InitRepository(path, bare)
		})

		Context("When path is invalid", func() {
			BeforeEach(func() {
				path = "/invalid"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When path is valid", func() {
			BeforeEach(func() {
				path = "testdata/git"
				resetDirectory(path)
			})

			AfterEach(func() {
				removeDirectory(path)
			})

			Context("When repository is not bare", func() {
				BeforeEach(func() {
					bare = false
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a valid Repository", func() {
					Expect(r).ToNot(BeNil())
					repo, ok := r.(*Repository)
					Expect(ok).To(BeTrue())
					dir, err := filepath.Abs(path)
					Expect(err).To(BeNil())
					Expect(repo.Path).To(Equal(dir))
				})
			})

			Context("With a bare repository", func() {
				BeforeEach(func() {
					bare = true
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a valid Repository", func() {
					Expect(r).ToNot(BeNil())
					repo, ok := r.(*Repository)
					Expect(ok).To(BeTrue())
					dir, err := filepath.Abs(path)
					Expect(err).To(BeNil())
					Expect(repo.Path).To(Equal(dir))
				})
			})
		})
	})

	Describe("OpenRepository", func() {
		var (
			r   vcs.Repository
			err error

			path string
		)

		JustBeforeEach(func() {
			r, err = OpenRepository(path)
		})

		Context("When path is invalid", func() {
			BeforeEach(func() {
				path = "/invalid"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When path is valid", func() {
			BeforeEach(func() {
				path = "testdata/git"
				resetDirectory(path)
			})

			AfterEach(func() {
				removeDirectory(path)
			})

			Context("When repository is not bare", func() {
				BeforeEach(func() {
					if o, err := exec.Command("git", "init", path).CombinedOutput(); err != nil {
						log.Fatalf("Could not init repository : %s", string(o))
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a valid Repository", func() {
					Expect(r).ToNot(BeNil())
					repo, ok := r.(*Repository)
					Expect(ok).To(BeTrue())
					dir, err := filepath.Abs(path)
					Expect(err).To(BeNil())
					Expect(repo.Path).To(Equal(dir))
				})

			})

			PContext("With a bare repository", func() {
				BeforeEach(func() {
					if o, err := exec.Command("git", "init", "--bare", path).CombinedOutput(); err != nil {
						log.Fatalf("Could not init repository : %s", string(o))
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a valid Repository", func() {
					Expect(r).ToNot(BeNil())
					repo, ok := r.(*Repository)
					Expect(ok).To(BeTrue())
					dir, err := filepath.Abs(path)
					Expect(err).To(BeNil())
					Expect(repo.Path).To(Equal(dir))
				})
			})
		})
	})

	Describe("Methods", func() {
		var (
			r *Repository
		)

		BeforeEach(func() {
			r = &Repository{
				Path: "testdata/repo",
			}
			resetDirectory(r.Path)
			if o, err := exec.Command("git", "init", r.Path).CombinedOutput(); err != nil {
				log.Fatalf("Could not init repository : %s", string(o))
			}
		})

		AfterEach(func() {
			removeDirectory(r.Path)
		})

		Describe("GetRevision", func() {
			var (
				revision string
				err      error
			)

			JustBeforeEach(func() {
				revision, err = r.GetRevision()
			})

			Context("With a new fresh git repository", func() {
				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With a valid git repository with some commits", func() {
				BeforeEach(func() {
					dir, _ := os.Getwd()
					os.Chdir(r.Path)
					defer os.Chdir(dir)
					if o, err := exec.Command("git", "commit", "--allow-empty", "-m", `"Test"`).CombinedOutput(); err != nil {
						log.Fatalf("Could not commit : %s", string(o))
					}
				})

				It("should return a valid revision", func() {
					Expect(revision).ToNot(BeZero())
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("GetRemote", func() {
			var (
				remote string
				err    error
			)

			JustBeforeEach(func() {
				remote, err = r.GetRemote()
			})

			Context("Without any remote", func() {
				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With a valid git repository with a remote", func() {
				var (
					remoteTest string
				)

				BeforeEach(func() {
					dir, _ := os.Getwd()
					os.Chdir(r.Path)
					defer os.Chdir(dir)
					remoteTest = "testremote"
					if o, err := exec.Command("git", "remote", "add", "origin", remoteTest).CombinedOutput(); err != nil {
						log.Fatalf("Could not add remote : %s", string(o))
					}
				})

				It("should return a valid revision", func() {
					Expect(remote).To(Equal(remoteTest))
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("AddRemote", func() {
			var (
				remote string
				err    error
			)

			JustBeforeEach(func() {
				err = r.AddRemote(remote)
			})

			Context("With an existing remote", func() {
				BeforeEach(func() {
					dir, _ := os.Getwd()
					os.Chdir(r.Path)
					defer os.Chdir(dir)
					remote = "testremote"
					if o, err := exec.Command("git", "remote", "add", "origin", remote).CombinedOutput(); err != nil {
						log.Fatalf("Could not add remote : %s", string(o))
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With a valid git repository without any remote", func() {
				BeforeEach(func() {
					remote = "test"
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		PDescribe("Fetch", func() {
		})

		Describe("Checkout", func() {
			var (
				err error

				revision string
			)

			JustBeforeEach(func() {
				err = r.Checkout(revision)
			})

			Context("With an unknown revision", func() {
				BeforeEach(func() {
					revision = "unknown"
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With a valid git repository without any remote", func() {
				BeforeEach(func() {
					dir, _ := os.Getwd()
					os.Chdir(r.Path)
					defer os.Chdir(dir)
					var rev []byte
					if o, err := exec.Command("git", "commit", "--allow-empty", "-m", `"Test"`).CombinedOutput(); err != nil {
						log.Fatalf("Could not commit : %s", string(o))
					} else if rev, err = exec.Command("git", "describe", "--always").CombinedOutput(); err != nil {
						log.Fatalf("Could not get revision : %s", string(o))
					}
					revision = string(rev[:len(rev)-1])
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Checkout", func() {
			var (
				err error

				branch string
			)

			JustBeforeEach(func() {
				err = r.CheckoutBranch(branch)
			})

			Context("With an unknown branch", func() {
				BeforeEach(func() {
					branch = "unknown"
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With a valid git repository without any remote", func() {
				BeforeEach(func() {
					dir, _ := os.Getwd()
					os.Chdir(r.Path)
					defer os.Chdir(dir)
					branch = "test"
					if o, err := exec.Command("git", "commit", "--allow-empty", "-m", `"Test"`).CombinedOutput(); err != nil {
						log.Fatalf("Could not commit : %s", string(o))
					} else if o, err := exec.Command("git", "branch", branch).CombinedOutput(); err != nil {
						log.Fatalf("Could not create branch : %s", string(o))
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("GetBranch", func() {
			var (
				branch string
				err    error
			)

			JustBeforeEach(func() {
				branch, err = r.GetBranch()
			})

			Context("Without any branch", func() {
				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With a valid branch", func() {
				BeforeEach(func() {
					dir, _ := os.Getwd()
					os.Chdir(r.Path)
					defer os.Chdir(dir)
					branch = "Branch"
					if o, err := exec.Command("git", "commit", "--allow-empty", "-m", `"Test"`).CombinedOutput(); err != nil {
						log.Fatalf("Could not commit : %s", string(o))
					} else if o, err := exec.Command("git", "checkout", "-b", branch).CombinedOutput(); err != nil {
						log.Fatalf("Could not commit : %s", string(o))
					}
				})

				It("should return a valid branch", func() {
					Expect(branch).To(Equal("Branch"))
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("GetPath", func() {
			var (
				path string

				err error
			)

			JustBeforeEach(func() {
				path, err = r.GetPath()
			})

			It("should return testdata/repo", func() {
				Expect(path).To(Equal("testdata/repo"))
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("GetVCSName", func() {
			var (
				name string
			)

			JustBeforeEach(func() {
				name = r.GetVCSName()
			})

			It("should return git", func() {
				Expect(name).To(Equal("git"))
			})
		})
	})
})
