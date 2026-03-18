package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

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

	// Windows GUI will load system info asynchronously to avoid UI freeze
	if runtime.GOOS == "windows" && uuidErr == nil {
		gui.RunWithSystemInfoAsync(uuidStr)
		return
	}

	gui.Run(uuidStr, uuidErr)
}
