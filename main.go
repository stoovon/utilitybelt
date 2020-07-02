package main

import (
	"github.com/stoovon/utilitybelt/cmd"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		panic(err)
	}
}
