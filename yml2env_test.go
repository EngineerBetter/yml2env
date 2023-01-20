package main_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"os/exec"
)

var cliPath string
var version string

var _ = BeforeSuite(func() {
	var err error
	data, err := os.ReadFile("version")
	version = string(data)
	Ω(err).ShouldNot(HaveOccurred())
	cliPath, err = Build("github.com/EngineerBetter/yml2env")
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	CleanupBuildArtifacts()
})

var _ = Describe("yml2env", func() {
	usage := "yml2env <YAML file> \\[<command> | --env\\]"

	It("requires a YAML file argument", func() {
		command := exec.Command(cliPath)
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Ω(session.Err).Should(Say(usage))
	})

	It("requires the YAML file to exist", func() {
		command := exec.Command(cliPath, "no/such/file.yml", "echo foo")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Ω(session.Err).Should(Say("no/such/file.yml does not exist"))
		Ω(session.Err).ShouldNot(Say("foo"))
	})

	Describe("running a command", func() {
		It("requires a command to invoke", func() {
			command := exec.Command(cliPath, "fixtures/vars.yml")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(1))
			Ω(session.Err).Should(Say(usage))
		})

		It("invokes the given command passing env vars from the YAML file", func() {
			command := exec.Command(cliPath, "fixtures/vars.yml", "fixtures/script.sh")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Ω(session).Should(Say("value from yaml"))
		})

		It("invokes the given command passing boolean env vars from the YAML file", func() {
			command := exec.Command(cliPath, "fixtures/boolean.yml", "fixtures/script.sh")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Ω(session).Should(Say("true"))
		})

		It("invokes the given command passing integer env vars from the YAML file", func() {
			command := exec.Command(cliPath, "fixtures/integer.yml", "fixtures/script.sh")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Ω(session).Should(Say("42"))
		})
	})

	Describe("printing out exports", func() {
		It("does not accept a subcommand", func() {
			command := exec.Command(cliPath, "fixtures/vars.yml", "--eval", "fixtures/script.sh")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(1))
			Ω(session.Err).Should(Say(usage))
		})

		It("prints out an export for each var", func() {
			command := exec.Command(cliPath, "fixtures/vars.yml", "--eval")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Ω(session).Should(Say("export 'VAR_FROM_YAML=value from yaml'"))
		})
	})

	Describe("checking the version", func() {
		It("returns the version in the version file when --version flag is provided", func() {
			command := exec.Command(cliPath, "--version")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Ω(session).Should(Say(version))
		})

		It("returns the version in the version file when -v flag is provided", func() {
			command := exec.Command(cliPath, "-v")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Ω(session).Should(Say(version))
		})
	})
})
