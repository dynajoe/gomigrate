package cmds

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/joeandaverde/gomigrate/core"
)

var InitDoc = `
Initializes a new project
`

var projectName string

func Init() error {
	config := core.NewConfig()

	dirs := []string{"deploy", "verify", "revert"}

	for _, p := range dirs {
		if err := os.MkdirAll(path.Join(config.RootDir, p), os.ModePerm); err != nil {
			return err
		}
	}

	if conf, err := os.Create("sqitch.conf"); err == nil {
		conf.WriteString(fmt.Sprintf(`# [core]\n    # engine =\n    # plan_file = %s.plan\n    # top_dir = .\n`, projectName))
		conf.Close()
	} else {
		return err
	}

	return core.MakePlanFile(config, projectName)
}

func InitFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&projectName, "project", "sqitch", "Project Name")
}
