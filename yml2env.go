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

var usage = "yml2env <YAML file> [<command> | --env]"

func main() {
	args := os.Args

	if isVersionCMD(args) {
		fmt.Fprintln(os.Stdout, version())
		os.Exit(0)
	}

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
	mapSlice = uppercaseKeys(mapSlice)
	envVars := os.Environ()
	envVars = addToEnv(mapSlice, envVars)

	if args[2] == "--eval" {
		if len(args) > 3 {
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(1)
		}

		printExports(mapSlice)
		os.Exit(0)
	} else {
		err, _ := run(envVars, args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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

func uppercaseKeys(mapSlice yaml.MapSlice) yaml.MapSlice {
	for i := 0; i < len(mapSlice); i++ {
		item := mapSlice[i]

		if key, ok := item.Key.(string); ok {
			key := strings.ToUpper(key)
			item = valueToString(item)
			if value, ok := item.Value.(string); ok {
				mapSlice[i] = yaml.MapItem{Key: key, Value: value}
			} else {
				fmt.Fprintln(os.Stderr, "YAML invalid")
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "YAML invalid")
			os.Exit(1)
		}
	}

	return mapSlice
}

func addToEnv(mapSlice yaml.MapSlice, envVars []string) []string {
	for i := 0; i < len(mapSlice); i++ {
		item := mapSlice[i]

		key, _ := item.Key.(string)
		item = valueToString(item)
		value, _ := item.Value.(string)
		envVars = env.Set(key, value, envVars)
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

func printExports(mapSlice yaml.MapSlice) {
	for i := 0; i < len(mapSlice); i++ {
		item := mapSlice[i]

		key, _ := item.Key.(string)
		key = strings.ToUpper(key)
		item = valueToString(item)
		value, _ := item.Value.(string)
		fmt.Printf("export '%s=%s'\n", key, value)
	}
}

func version() string {
	data, err := os.ReadFile("version")
	if err != nil {
		fmt.Printf("unable to retrieve version: %s", err)
	}
	return string(data)
}

func isVersionCMD(args []string) bool {
	if len(args) > 1 && (args[1] == "--version" || args[1] == "-v") {
		return true
	}
	return false
}
