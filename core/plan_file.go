package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// PlanFile represents a migration plan
type PlanFile struct {
	Project    string
	Name       string
	Path       string
	Migrations []Migration
	f          *os.File
}

// MakePlanFile creates a new plan file. If one exists an error is returned.
func MakePlanFile(config Config, name string) error {
	planPath := filepath.Join(config.RootDir, name+".plan")

	header := []string{"%syntax-version=1.0.0", "%project=projectname", "%uri=https://project"}

	if _, err := os.Stat(planPath); os.IsNotExist(err) {
		planFile, err := os.OpenFile(planPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

		if err != nil {
			panic(err)
		}

		defer planFile.Close()

		planFile.WriteString(strings.Join(header, "\n") + "\n")
	}

	return nil
}

type void struct{}

var member void

// LoadPlan creates a new plan file. The file is expected to exist.
func LoadPlan(config Config, name string) *PlanFile {
	planPath := filepath.Join(config.RootDir, name+".plan")

	data, err := ioutil.ReadFile(planPath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	var migrations []Migration
	var migrationKeys = make(map[string]void)
	var tags []string

	for _, line := range lines {
		firstRune, _ := utf8.DecodeRuneInString(line)

		if len(line) > 0 && unicode.IsLetter(firstRune) {
			if migration, err := ParseMigration(line); err == nil {
				if _, exists := migrationKeys[migration.Name]; exists {
					panic("Duplicate migration")
				} else {
					migrationKeys[migration.Name] = member
				}

				content, _ := ioutil.ReadFile(path.Join(config.RootDir, "deploy", migration.Name+".sql"))
				migration.Content = string(content)

				migrations = append(migrations, migration)
			} else {
				panic(err)
			}
		} else if firstRune == rune('@') {
			tags = append(tags, line)
		}
	}

	f, err := os.OpenFile(planPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return &PlanFile{
		Project:    "app",
		Name:       name,
		Path:       planPath,
		Migrations: migrations,
		f:          f,
	}
}

func migrationTemplate(plan string, name string) string {
	return fmt.Sprintf("-- Deploy %s:%s to pg\nBEGIN;\n\n-- DDL HERE\n\nCOMMIT;", plan, name)
}

// AddMigration adds a migration to the end of the plan file
func (plan *PlanFile) AddMigration(config Config, name string, comment string) error {
	now := time.Now().UTC()

	migrationPath := filepath.Join(config.RootDir, "deploy", name+".sql")

	dirName := filepath.Dir(migrationPath)

	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return err
	}

	migrationFile, err := os.Create(migrationPath)

	if err != nil {
		return err
	}

	defer migrationFile.Close()

	if _, err := migrationFile.WriteString(migrationTemplate(plan.Name, name)); err != nil {
		return err
	}

	if err := migrationFile.Sync(); err != nil {
		return err
	}

	if _, err := plan.f.WriteString(fmt.Sprintf("%s %s %s <%s> # %s\n",
		name,
		now.Format("2006-01-02T15:04:05Z"),
		config.Name,
		config.Email,
		comment)); err != nil {
		panic(err)
	}

	return nil
}

// Close closes the plan file
func (plan *PlanFile) Close() {
	plan.f.Sync()
	plan.f.Close()
}
