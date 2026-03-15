package main

import (
	"flag"
	"fmt"
	"os"

	"uuidgen/gui"
	"uuidgen/uuid"
)

func main() {
	cliMode := flag.Bool("cli", false, "Run in CLI mode (no GUI)")
	flag.Parse()

	uuidStr, uuidErr := uuid.GetUUID()

	if *cliMode {
		if uuidErr != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", uuidErr)
			os.Exit(1)
		}
		fmt.Println(uuidStr)
		return
	}

	gui.Run(uuidStr, uuidErr)
}
