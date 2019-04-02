package cmds

import (
	"flag"

	"github.com/joeandaverde/gomigrate/core"
)

var AddChangeDoc = `
Add a new change
`

var changeName string
var planName string

func AddChange() error {
	config := core.NewConfig()
	plan := core.LoadPlan(config, planName)

	return plan.AddChange(config, changeName, "Note")
}

func AddChangeFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&changeName, "name", "", "change name")
	fs.StringVar(&planName, "plan", "sqitch", "plan name")
}
