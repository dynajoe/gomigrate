package cmds

import (
	"crypto/sha1"
	"database/sql"
	"flag"
	"fmt"
	"unicode/utf8"

	"github.com/joeandaverde/gomigrate/core"
	"github.com/joeandaverde/gomigrate/data"
)

var DeployDoc = `
Specify database using --target flag.
`

var targetDatabase string

func ensureRegistry(db *sql.DB) {
	rows, err := db.Query("SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'sqitch';")

	if err != nil {
		panic(err)
	}

	if rows.Next() {
		return
	}

	fmt.Println("Deploying registry...")
	registrySQL := data.MakeSqitchRegistrySQL("sqitch")
	_, err = db.Query(registrySQL)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func deployChange(db *sql.DB, config core.Config, plan *core.PlanFile, change core.Change) {
	scriptHash := fmt.Sprintf("%x", sha1.Sum([]byte(change.Content)))
	contentUTF8 := runesToUTF8Manual([]rune(change.Content))
	changeID := fmt.Sprintf("%x", sha1.Sum([]byte("change "+string(len(contentUTF8))+"\000"+string(contentUTF8))))

	rows, err := db.Query("SELECT * FROM sqitch.changes WHERE project = $1 AND script_hash = $2;",
		plan.Project, scriptHash)

	if rows.Next() {
		fmt.Println("Change already deployed")
		return
	}

	fmt.Printf("+ %s .. ", change.Name)

	_, err = db.Exec(change.Content)

	if err == nil {
		fmt.Println("ok")

		_, err := db.Exec(`
			INSERT INTO sqitch.changes (change_id, script_hash, change, project, note, committer_name, committer_email, planned_at, planner_name, planner_email)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
			changeID, scriptHash, change.Name, plan.Project, change.Comment, config.User, config.Email, change.Date, change.CreatedBy, change.CreatedBy)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func runesToUTF8Manual(rs []rune) []byte {
	bs := make([]byte, len(rs)*utf8.UTFMax)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs[:count]
}

func Deploy() error {
	config := core.NewConfig()

	plan := core.LoadPlan(config, planName)

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	ensureRegistry(db)

	for _, m := range plan.Changes {
		deployChange(db, config, plan, m)
	}

	return nil
}

func DeployFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&targetDatabase, "url", "", "Target database")
	fs.StringVar(&planName, "plan", "sqitch", "Plan")
}
