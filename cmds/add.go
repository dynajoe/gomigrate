package cmds

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var AddMigrationDoc = `
Add a new migration
`

var migrationName string

func AddMigration() error {
	migrationPath := filepath.Join("migrations", migrationName+".sql")
	planPath := filepath.Join("migrations", "app.plan")

	dirName := filepath.Dir(migrationPath)
	err := os.MkdirAll(dirName, os.ModePerm)
	migrationFile, err := os.Create(migrationPath)
	defer migrationFile.Close()

	if err != nil {
		return err
	}

	migrationFile.WriteString("BEGIN;\n\n-- DDL HERE\n\nCOMMIT;")
	migrationFile.Sync()

	planFile, err := os.OpenFile(planPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer planFile.Close()

	t := time.Now()
	planFile.WriteString(fmt.Sprintf("\n%s %s", migrationName, t.Format(time.RFC3339)))

	return nil
}

func AddMigrationFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&migrationName, "name", "", "migration name")
}
