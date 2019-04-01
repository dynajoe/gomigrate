package cmds

import (
	"crypto/sha1"
	"database/sql"
	"flag"
	"fmt"

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
	fmt.Println(registrySQL)
	_, err = db.Query(registrySQL)

	if err != nil {
		panic(err)
	}
}

func deployChange(db *sql.DB, config core.Config, plan *core.PlanFile, change core.Migration) {
	scriptHash := sha1.New()
	scriptHash.Write([]byte(change.Content))
	hashValue := fmt.Sprintf("%x", scriptHash.Sum(nil))
	var changeExists bool

	err := db.QueryRow("SELECT true FROM sqitch.changes WHERE project = $1 AND script_hash = $2;", plan.Project, hashValue).Scan(&changeExists)

	if err != nil && err != sql.ErrNoRows {
		panic("Error checking for migration")
	}

	if changeExists {
		fmt.Println("Change already deployed")
		return
	}

	fmt.Printf("+ %s .. ", change.Name)

	_, err = db.Exec(change.Content)

	if err != nil {
		fmt.Println("ok")
		_, err := db.Exec(`
			INSERT INTO sqitch.changes(script_hash, change, project, note, committer_name, comitter_email, planned_at, planner_name, planner_email)
			VALUES ($1, $2, $3, $4, $6, $7, $8, $9, $10);`, hashValue, change.Name, plan.Project, change.Comment, config.User, config.Email, change.Date, change.CreatedBy, change.CreatedBy)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func Deploy() error {
	config := core.NewConfig()

	plan := core.LoadPlan(config, planName)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	ensureRegistry(db)

	for _, m := range plan.Migrations {
		deployChange(db, config, plan, m)
	}

	return nil
}

func DeployFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&targetDatabase, "url", "", "Target database")
	fs.StringVar(&planName, "plan", "sqitch", "Plan")
}
