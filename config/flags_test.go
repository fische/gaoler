package config_test

import (
	"os"

	. "github.com/fische/gaoler/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flags", func() {
	var (
		f Flags
	)

	Describe("OpenFlags", func() {
		var (
			ret int
		)

		JustBeforeEach(func() {
			ret = f.OpenFlags()
		})

		Context("With Save flag", func() {
			BeforeEach(func() {
				f = Save
			})

			It("should return os.O_WRONLY|os.O_CREATE", func() {
				Expect(ret).To(Equal(os.O_WRONLY | os.O_CREATE))
			})
		})

		Context("With Load flag", func() {
			BeforeEach(func() {
				f = Load
			})

			It("should return os.O_RDONLY", func() {
				Expect(ret).To(Equal(os.O_RDONLY))
			})
		})

		Context("With Load and Save flags", func() {
			BeforeEach(func() {
				f = Load | Save
			})

			It("should return os.O_RDWR|os.O_CREATE", func() {
				Expect(ret).To(Equal(os.O_RDWR | os.O_CREATE))
			})
		})

		Context("Without any known flags", func() {
			BeforeEach(func() {
				f = 0
			})

			It("should return 0", func() {
				Expect(ret).To(Equal(0))
			})
		})
	})

	Describe("Has", func() {
		var (
			o Flags

			ret bool
		)

		JustBeforeEach(func() {
			ret = f.Has(o)
		})

		Context("Without flag", func() {
			BeforeEach(func() {
				f = Load
				o = Save
			})

			It("shoult return false", func() {
				Expect(ret).To(BeFalse())
			})
		})

		Context("With flag", func() {
			BeforeEach(func() {
				f = Load | Save
				o = Load
			})

			It("shoult return true", func() {
				Expect(ret).To(BeTrue())
			})
		})
	})
})
