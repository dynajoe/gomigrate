package cmds

import (
	"flag"

	"github.com/joeandaverde/gomigrate/core"
)

var AddMigrationDoc = `
Add a new migration
`

var migrationName string
var planName string

func AddMigration() error {
	config := core.NewConfig()
	plan := core.LoadPlan(config, planName)

	return plan.AddMigration(config, migrationName, "Note")
}

func AddMigrationFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&migrationName, "name", "", "migration name")
	fs.StringVar(&planName, "plan", "sqitch", "plan name")
}
