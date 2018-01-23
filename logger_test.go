package libbuildpack_test

import (
	"bytes"

	"github.com/cloudfoundry/libbuildpack"
	env "github.com/cloudfoundry/libbuildpack/env"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {

	var (
		logger  *libbuildpack.Logger
		buffer  *bytes.Buffer
		mockEnv env.Env
	)

	BeforeEach(func() {
		buffer = new(bytes.Buffer)
		mockEnv = env.Mock()
		logger = libbuildpack.NewLogger(buffer, mockEnv)
	})

	Describe("Debug", func() {
		Context("BP_DEBUG is set", func() {
			BeforeEach(func() {
				Expect(mockEnv.Set("BP_DEBUG", "true")).To(Succeed())
			})

			It("Logs the message", func() {
				logger.Debug("detailed info")
				Expect(buffer.String()).To(ContainSubstring("\033[34;1mDEBUG:\033[0m detailed info"))
			})
		})

		Context("BP_DEBUG is not set", func() {
			BeforeEach(func() {
				Expect(mockEnv.Set("BP_DEBUG", "")).To(Succeed())
			})

			It("Does not log the message", func() {
				logger.Debug("detailed info")
				Expect(buffer.String()).To(Equal(""))
			})
		})
	})
})
