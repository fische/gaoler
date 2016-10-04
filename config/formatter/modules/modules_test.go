package modules_test

import (
	"github.com/fische/gaoler/config/formatter"
	. "github.com/fische/gaoler/config/formatter/modules"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Modules", func() {
	Describe("Register / GetFormatter / Formatters", func() {
		var (
			v          *factory
			ret        formatter.Factory
			registered []string
		)

		JustBeforeEach(func() {
			Register(v)
			ret = GetFormatter(v.types[0])
			registered = Formatters()
		})

		Context("With a new Formatter", func() {
			BeforeEach(func() {
				v = &factory{
					types: []string{"test"},
				}
			})

			It("should have correctly registered Formatter", func() {
				Expect(ret).To(Equal(v))
				passed := false
				for _, reg := range registered {
					if reg == v.types[0] {
						passed = true
					}
				}
				Expect(passed).To(BeTrue())
			})
		})
	})

	Describe("GetFormatter", func() {
		var (
			ext string

			ret formatter.Factory
		)

		JustBeforeEach(func() {
			ret = GetFormatter(ext)
		})

		Context("With an unknown Formatter", func() {
			BeforeEach(func() {
				ext = "___unknown___"
			})

			It("should return nil", func() {
				Expect(ret).To(BeNil())
			})
		})
	})
})
