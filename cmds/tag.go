package cmds

// @v1.0.0 2018-09-28T19:00:24Z Joe Andaverde <joe.andaverde@smrxt.com> # v1.0.0

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
