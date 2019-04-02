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

	addChangeCmd := comandante.NewCommand("add", "add a new migration", cmds.AddChange)
	addChangeCmd.FlagInit = cmds.AddChangeFlagHandler
	addChangeCmd.Documentation = cmds.AddChangeDoc
	bin.RegisterCommand(addChangeCmd)

	deployCmd := comandante.NewCommand("deploy", "deploy migrations", cmds.Deploy)
	deployCmd.FlagInit = cmds.DeployFlagHandler
	deployCmd.Documentation = cmds.DeployDoc
	bin.RegisterCommand(deployCmd)

	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
