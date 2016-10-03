package vcs_test

import (
	"errors"

	. "github.com/fische/gaoler/vcs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VCS", func() {
	var (
		r   Repository
		err error

		v    *VCSTest
		path string

		ret *RepositoryTest
	)

	BeforeEach(func() {
		v = &VCSTest{}
		ret = &RepositoryTest{
			name: "test",
		}
	})

	Describe("CloneAtRevision", func() {
		var (
			remote   string
			revision string
		)

		BeforeEach(func() {
			remote = "remote"
			revision = "revision"
		})

		JustBeforeEach(func() {
			r, err = CloneAtRevision(v, remote, revision, path)
		})

		Context("When repository creation fails", func() {
			BeforeEach(func() {
				v.new = func(path string) (Repository, error) {
					return nil, errors.New("")
				}
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When repository creation succeeds", func() {
			BeforeEach(func() {
				v.new = func(path string) (Repository, error) {
					return ret, nil
				}
			})

			Context("When adding remote fails", func() {
				BeforeEach(func() {
					ret.addRemote = func(remoteTest string) error {
						return errors.New("")
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("When adding remote succeeds", func() {
				BeforeEach(func() {
					ret.addRemote = func(remoteTest string) error {
						Expect(remoteTest).To(Equal(remote))
						return nil
					}
				})

				Context("When fetching fails", func() {
					BeforeEach(func() {
						ret.fetch = func() error {
							return errors.New("")
						}
					})

					It("should return an error", func() {
						Expect(err).ToNot(BeNil())
					})
				})

				Context("When fetching succeeds", func() {
					var (
						fetched bool
					)

					BeforeEach(func() {
						fetched = false
						ret.fetch = func() error {
							fetched = true
							return nil
						}
					})

					Context("When checking out fails", func() {
						BeforeEach(func() {
							ret.checkout = func(rev string) error {
								return errors.New("")
							}
						})

						It("should return an error", func() {
							Expect(err).ToNot(BeNil())
						})
					})

					Context("When checking out fails", func() {
						BeforeEach(func() {
							ret.checkout = func(rev string) error {
								Expect(rev).To(Equal(revision))
								return nil
							}
						})

						It("should not return an error", func() {
							Expect(err).To(BeNil())
						})

						It("should return a valid repository", func() {
							Expect(r).To(Equal(ret))
						})

						It("should have fetched", func() {
							Expect(fetched).To(BeTrue())
						})
					})
				})
			})
		})
	})

	Describe("CloneRepository", func() {
		var (
			repo *RepositoryTest
		)

		BeforeEach(func() {
			repo = &RepositoryTest{}
			path = "/valid"
		})

		JustBeforeEach(func() {
			r, err = CloneRepository(v, repo, path)
		})

		Context("When getting revision fails", func() {
			BeforeEach(func() {
				repo.revision = func() (string, error) {
					return "", errors.New("")
				}
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When getting revision succeeds", func() {
			var (
				revision string
			)

			BeforeEach(func() {
				revision = "Revision"
				repo.revision = func() (string, error) {
					return revision, nil
				}
			})

			Context("When getting remote fails", func() {
				BeforeEach(func() {
					repo.remote = func() (string, error) {
						return "", errors.New("")
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("When getting remote succeeds", func() {
				var (
					remote string
				)

				BeforeEach(func() {
					remote = "Remote"
					repo.remote = func() (string, error) {
						return remote, nil
					}
				})

				Context("When cloning fails", func() {
					BeforeEach(func() {
						v.new = func(path string) (Repository, error) {
							return nil, errors.New("")
						}
					})

					It("should return an error", func() {
						Expect(err).ToNot(BeNil())
					})
				})

				Context("When cloning succeeds", func() {
					var (
						fetched bool
					)

					BeforeEach(func() {
						fetched = false
						v.new = func(path string) (Repository, error) {
							return ret, nil
						}
						ret.addRemote = func(remoteTest string) error {
							Expect(remoteTest).To(Equal(remote))
							return nil
						}
						ret.fetch = func() error {
							fetched = true
							return nil
						}
						ret.checkout = func(rev string) error {
							Expect(rev).To(Equal(revision))
							return nil
						}
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should return a valid repository", func() {
						Expect(r).To(Equal(ret))
					})

					It("should have fetched", func() {
						Expect(fetched).To(BeTrue())
					})
				})
			})
		})
	})
})
