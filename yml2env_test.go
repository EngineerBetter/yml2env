package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"os/exec"
)

var _ = Describe("yml2env", func() {
	var cliPath string
	usage := "yml2env <YAML file> \\[<command> | --env\\]"

	BeforeSuite(func() {
		var err error
		cliPath, err = Build("github.com/EngineerBetter/yml2env")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		CleanupBuildArtifacts()
	})

	It("requires a YAML file argument", func() {
		command := exec.Command(cliPath)
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Ω(session.Err).Should(Say(usage))
	})

	It("requires a the YAML file to exist", func() {
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
})
