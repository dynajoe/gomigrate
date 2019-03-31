package main

import (
	"fmt"

	"os"

	"github.com/joeandaverde/gomigrate/cmds"
	"github.com/robmerrell/comandante"
)

func main() {
	bin := comandante.New("gomigrate", "Example program showing how to use Comandante")
	bin.IncludeHelp()

	addMigrationCmd := comandante.NewCommand("add", "add a new migration", cmds.AddMigration)
	addMigrationCmd.FlagInit = cmds.AddMigrationFlagHandler
	addMigrationCmd.Documentation = cmds.AddMigrationDoc
	bin.RegisterCommand(addMigrationCmd)

	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
