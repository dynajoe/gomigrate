package main

import (
	"fmt"

	"os"

	"github.com/joeandaverde/gomigrate/cmds"
	_ "github.com/lib/pq"
	"github.com/robmerrell/comandante"
)

func main() {
	bin := comandante.New("gomigrate", "Example program showing how to use Comandante")
	bin.IncludeHelp()

	addMigrationCmd := comandante.NewCommand("add", "add a new migration", cmds.AddMigration)
	addMigrationCmd.FlagInit = cmds.AddMigrationFlagHandler
	addMigrationCmd.Documentation = cmds.AddMigrationDoc
	bin.RegisterCommand(addMigrationCmd)

	deployCmd := comandante.NewCommand("deploy", "deploy migrations", cmds.Deploy)
	deployCmd.FlagInit = cmds.DeployFlagHandler
	deployCmd.Documentation = cmds.DeployDoc
	bin.RegisterCommand(deployCmd)

	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
