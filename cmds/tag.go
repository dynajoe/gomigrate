package cmds

import (
	"flag"
)

var TagDoc = `
Specify tag using --tag flag.
`

var planTag string

func Tag() error {
	return nil
}

func TagFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&planTag, "tag", "", "Tag")
}
