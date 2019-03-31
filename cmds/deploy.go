package cmds

import (
	"flag"
)

var DeployDoc = `
Specify database using --target flag.
`

var targetDatabase string

func Deploy() error {
	return nil
}

func DeployFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&targetDatabase, "url", "", "Target database")
}
