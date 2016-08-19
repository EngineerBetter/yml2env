package main

import (
	"bytes"
	"fmt"
	"github.com/EngineerBetter/yml2env/env"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var usage = "yml2env <YAML file> <command>"

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	yamlPath := args[1]

	if !fileExists(yamlPath) {
		fmt.Fprintln(os.Stderr, yamlPath+" does not exist")
		os.Exit(1)
	}

	//Load YAML
	bytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read "+yamlPath)
	}

	//Get as map
	vars := yaml.MapSlice{}
	err = yaml.Unmarshal([]byte(bytes), &vars)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse "+yamlPath)
		os.Exit(1)
	}

	envVars := os.Environ()

	//uppercase keys
	for i := 0; i < len(vars); i++ {
		item := vars[i]

		if key, ok := item.Key.(string); ok {
			key := strings.ToUpper(key)
			if value, ok := item.Value.(string); ok {
				env.Set(key, value, envVars)
				fmt.Println(key, value)
			} else {
				fmt.Fprintln(os.Stderr, "YAML invalid")
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "YAML invalid")
			os.Exit(1)
		}
	}

	run(envVars, args[1:])
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func commandWithEnv(envVars []string, args ...string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = envVars
	return cmd
}

func run(envVars []string, args []string) (error, int, string) {
	cmd := commandWithEnv(envVars, args...)

	buffer := bytes.NewBufferString("")
	multiWriter := io.MultiWriter(os.Stdout, buffer)

	cmd.Stdin = os.Stdin
	cmd.Stdout = multiWriter
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		return err, -1, ""
	}

	err = cmd.Wait()
	output := buffer.String()
	return nil, determineExitCode(cmd, err), output
}

func determineExitCode(cmd *exec.Cmd, err error) (exitCode int) {
	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if status.Signaled() {
		exitCode = 128 + int(status.Signal())
	} else {
		exitStatus := status.ExitStatus()
		if exitStatus == -1 && err != nil {
			exitCode = 254
		}
		exitCode = exitStatus
	}

	return
}
