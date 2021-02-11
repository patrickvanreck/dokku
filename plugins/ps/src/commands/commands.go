package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

const (
	helpHeader = `Usage: dokku ps[:COMMAND]

Manage app processes

Additional commands:`

	helpContent = `
    ps:inspect <app>, Displays a sanitized version of docker inspect for an app
    ps:rebuild [--parallel count] [--all|<app>], Rebuilds an app from source
    ps:report [<app>] [<flag>], Displays a process report for one or more apps
    ps:restart [--parallel count] [--all|<app>], Restart an app
    ps:restore [<app>], Start previously running apps e.g. after reboot
    ps:scale [--skip-deploy] <app> <proc>=<count> [<proc>=<count>...], Get/Set how many instances of a given process to run
    ps:set <app> <key> <value>, Set or clear a ps property for an app
    ps:start [--parallel count] [--all|<app>], Start an app
    ps:stop [--parallel count] [--all|<app>], Stop an app
`
)

func main() {
	flag.Usage = usage
	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case "ps", "ps:help":
		usage()
	case "help":
		command := common.NewShellCmd(fmt.Sprintf("ps -o command= %d", os.Getppid()))
		command.ShowOutput = false
		output, err := command.Output()

		if err == nil && strings.Contains(string(output), "--all") {
			fmt.Println(helpContent)
		} else {
			fmt.Print("\n    ps, Manage app processes\n")
		}
	default:
		dokkuNotImplementExitCode, err := strconv.Atoi(os.Getenv("DOKKU_NOT_IMPLEMENTED_EXIT"))
		if err != nil {
			fmt.Println("failed to retrieve DOKKU_NOT_IMPLEMENTED_EXIT environment variable")
			dokkuNotImplementExitCode = 10
		}
		os.Exit(dokkuNotImplementExitCode)
	}
}

func usage() {
	common.CommandUsage(helpHeader, helpContent)
}
