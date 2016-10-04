package config_test

import (
	"errors"
	"log"
	"os"

	. "github.com/fische/gaoler/config"
	"github.com/fische/gaoler/config/formatter/modules"
	"github.com/fische/gaoler/project"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		err error
		cfg *Config
	)

	Describe("New", func() {
		var (
			p          *project.Project
			configPath string
			flags      Flags
		)

		BeforeEach(func() {
			p = &project.Project{}
		})

		JustBeforeEach(func() {
			cfg, err = New(p, configPath, flags)
		})

		Context("When opening file fails", func() {
			BeforeEach(func() {
				configPath = "/nowhere"
				flags = Load
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When formatter is not found", func() {
			BeforeEach(func() {
				configPath = "testdata/test.unknown"
				flags = Load
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When creating config succeeds", func() {
			var (
				f *factory
			)

			BeforeEach(func() {
				configPath = "testdata/test.testfile"
				flags = Load
				f = &factory{
					types: []string{"testfile"},
				}
				modules.Register(f)
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should set Config correctly", func() {
				Expect(cfg.File).ToNot(BeNil())
				Expect(cfg.File.Name()).To(Equal("testdata/test.testfile"))
				Expect(cfg.Format).To(Equal(f))
				Expect(cfg.Project).To(Equal(p))
			})
		})
	})

	Describe("Methods", func() {
		BeforeEach(func() {
			cfg = &Config{
				Project: &project.Project{
					Name: "test",
				},
			}
		})

		Describe("Save", func() {
			JustBeforeEach(func() {
				err = cfg.Save()
			})

			Context("When file is a directory", func() {
				BeforeEach(func() {
					cfg.File, err = os.Open("testdata")
					if err != nil {
						log.Fatalf("Could not open testdata : %v", err)
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("When file is nil", func() {
				BeforeEach(func() {
					cfg.File = (*os.File)(nil)
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With a valid file", func() {
				BeforeEach(func() {
					cfg.File, err = os.OpenFile("testdata/test.testfile", os.O_RDWR, 0664)
					if err != nil {
						log.Fatalf("Could not open testdata : %v", err)
					}
				})

				Context("With a valid Encoder", func() {
					BeforeEach(func() {
						cfg.Format = &factory{
							encoder: encoder{
								encode: func(i interface{}) error {
									Expect(i).To(Equal(cfg.Project))
									return nil
								},
							},
						}
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})
				})

				Context("With an invalid Encoder", func() {
					BeforeEach(func() {
						cfg.Format = &factory{
							encoder: encoder{
								encode: func(i interface{}) error {
									return errors.New("")
								},
							},
						}
					})

					It("should return an error", func() {
						Expect(err).ToNot(BeNil())
					})
				})

				Context("With a valid PrettyEncoder", func() {
					BeforeEach(func() {
						cfg.Format = &factory{
							encoder: prettyEncoder{
								prettyEncode: func(i interface{}) error {
									Expect(i).To(Equal(cfg.Project))
									return nil
								},
							},
						}
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})
				})

				Context("With an invalid PrettyEncoder", func() {
					BeforeEach(func() {
						cfg.Format = &factory{
							encoder: prettyEncoder{
								prettyEncode: func(i interface{}) error {
									return errors.New("")
								},
							},
						}
					})

					It("should return an error", func() {
						Expect(err).ToNot(BeNil())
					})
				})
			})
		})

		Describe("Load", func() {
			JustBeforeEach(func() {
				err = cfg.Load()
			})

			Context("When decoding fails", func() {
				BeforeEach(func() {
					cfg.Format = &factory{
						decoder: decoder{
							decode: func(i interface{}) error {
								return errors.New("")
							},
						},
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("When decoding succeed", func() {
				BeforeEach(func() {
					cfg.Format = &factory{
						decoder: decoder{
							decode: func(i interface{}) error {
								Expect(i).To(Equal(cfg.Project))
								return nil
							},
						},
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
