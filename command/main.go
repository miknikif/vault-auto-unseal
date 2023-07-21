package command

import "fmt"

// Main Run function
func Run(args []string) int {
	err := StartHttpServer()

	if err != nil {
		fmt.Println(fmt.Errorf("%s", err))
		return 1
	}

	return 0
}
