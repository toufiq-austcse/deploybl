package cmd_runner

import (
	"fmt"
	"os/exec"
)

func RunCommand(name string, args []string) (*string, error) {
	cmd := exec.Command(name, args...)
	fmt.Println("executing ", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return nil, err
	}

	outputString := string(output)
	fmt.Println("output: ", outputString)

	return &outputString, nil
}