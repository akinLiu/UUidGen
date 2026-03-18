package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"uuidgen/gui"
	"uuidgen/sysinfo"
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

	// Get full system info on Windows
	if runtime.GOOS == "windows" && uuidErr == nil {
		info, err := sysinfo.GetSystemInfo(uuidStr)
		if err == nil {
			gui.RunWithSystemInfo(info)
			return
		}
	}

	gui.Run(uuidStr, uuidErr)
}
