package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/EngineerBetter/yml2env/env"
	"gopkg.in/yaml.v2"
)

var usage = "yml2env <YAML file> <command>"

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	yamlPath := args[1]

	if !fileExists(yamlPath) {
		fmt.Fprintln(os.Stderr, yamlPath+" does not exist")
		os.Exit(1)
	}

	bytes := loadYaml(yamlPath)
	mapSlice := parseYaml(bytes)
	envVars := os.Environ()
	envVars = addUppercaseKeysToEnv(mapSlice, envVars)

	err, _ := run(envVars, args[2:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func loadYaml(yamlPath string) []byte {
	bytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read "+yamlPath)
	}
	return bytes
}

func parseYaml(bytes []byte) yaml.MapSlice {
	vars := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(bytes), &vars)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse YAML")
		os.Exit(1)
	}

	return vars
}

func valueToString(item yaml.MapItem) yaml.MapItem {
	if value, ok := item.Value.(bool); ok {
		item.Value = strconv.FormatBool(value)
	} else if value, ok := item.Value.(int); ok {
		item.Value = strconv.Itoa(value)
	}
	return item
}

func addUppercaseKeysToEnv(mapSlice yaml.MapSlice, envVars []string) []string {
	for i := 0; i < len(mapSlice); i++ {
		item := mapSlice[i]

		if key, ok := item.Key.(string); ok {
			key := strings.ToUpper(key)
			item = valueToString(item)
			if value, ok := item.Value.(string); ok {
				envVars = env.Set(key, value, envVars)
			} else {
				fmt.Fprintln(os.Stderr, "YAML invalid")
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "YAML invalid")
			os.Exit(1)
		}
	}

	return envVars
}

func commandWithEnv(envVars []string, args ...string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = envVars
	return cmd
}

func run(envVars []string, args []string) (error, int) {
	cmd := commandWithEnv(envVars, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		return err, -1
	}

	err = cmd.Wait()
	return nil, determineExitCode(cmd, err)
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
