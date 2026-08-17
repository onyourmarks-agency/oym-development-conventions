package main

import (
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"

	"github.com/onyourmarks-agency/oym-development-conventions/cmd/conventions"
)

//go:embed conventions templates
var assets embed.FS

func main() {
	err := conventions.New(assets).Execute()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
