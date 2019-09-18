package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Too few arguments")
	}

	dir := os.Args[1]

	command := os.Args[2]
	restArgs := os.Args[3:]

	envVars, err := getEnvVars(dir)
	if err != nil {
		log.Fatal(err)
	}

	stat, _ := os.Stdin.Stat()
	isStdinPiped := (stat.Mode() & os.ModeCharDevice) > 0

	err = execCommand(isStdinPiped, envVars, command, restArgs)
	if err != nil {
		log.Fatalf("Error executing program: %v", err)
	}
}

func getEnvVars(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var vars []string

	for _, f := range files {
		contents, err := ioutil.ReadFile(dir + "/" + f.Name())
		if err != nil {
			return nil, err
		}
		vars = append(vars, f.Name()+"="+string(contents))
	}
	return vars, nil
}

func execCommand(isStdinPiped bool, envVars []string, command string, restArgs []string) error {
	cmd := exec.Command(command, restArgs...)
	cmd.Stdout = os.Stdout
	if isStdinPiped {
		cmd.Stdin = os.Stdin
	}
	cmd.Env = append(os.Environ(), envVars...)

	return cmd.Run()
}
