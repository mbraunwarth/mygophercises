package main

import "github.com/spf13/cobra/cobra/cmd"

/* CLI Task Manager
 * This will be a command line app to manage tasks. The app should include
 * at least the functionality to `add`, `delete` and `list` tasks.
 *
 * Usage: task <command> [<args>]
 */

func main() {
	cmd.Execute()
}
