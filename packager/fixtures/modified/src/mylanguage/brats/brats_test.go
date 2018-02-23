package brats_test

import (
	"github.com/cloudfoundry/libbuildpack/bratshelper"
	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO explain when to make not pending
var _ = PDescribe("mylanguage buildpack", func() {
	bratshelper.UnbuiltBuildpack("mylanguage", CopyBrats)
	bratshelper.DeployingAnAppWithAnUpdatedVersionOfTheSameBuildpack(CopyBrats)
	bratshelper.StagingWithBuildpackThatSetsEOL("mylanguage", CopyBrats)
	bratshelper.StagingWithADepThatIsNotTheLatest("mylanguage", CopyBrats)
	bratshelper.StagingWithCustomBuildpackWithCredentialsInDependencies(`mylanguage\-[\d\.]+\-linux\-x64\-[\da-f]+\.tgz`, CopyBrats)
	bratshelper.DeployAppWithExecutableProfileScript("mylanguage", CopyBrats)
	bratshelper.DeployAnAppWithSensitiveEnvironmentVariables(CopyBrats)
	bratshelper.ForAllSupportedVersions("mylanguage", CopyBrats, func(version string, app *cutlass.App) {
		PushApp(app)

		By("does a thing", func() {
			Expect(app).ToNot(BeNil())
		})
	})
})
