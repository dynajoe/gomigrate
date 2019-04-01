package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PlanFile represents a migration plan
type PlanFile struct {
	Name string
	Path string
	f    *os.File
}

// MakePlanFile creates a new plan file. If one exists an error is returned.
func MakePlanFile(config *Config, name string) error {
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

// LoadPlan creates a new plan file. The file is expected to exist.
func LoadPlan(config *Config, name string) *PlanFile {
	planPath := filepath.Join(config.RootDir, name+".plan")

	data, err := ioutil.ReadFile(planPath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		fmt.Println(line)
	}

	f, err := os.OpenFile(planPath, os.O_APPEND|os.O_WRONLY, 0644)
	return &PlanFile{
		Name: name,
		Path: planPath,
		f:    f,
	}
}

func migrationTemplate(plan string, name string) string {
	return fmt.Sprintf("-- Deploy %s:%s to pg\nBEGIN;\n\n-- DDL HERE\n\nCOMMIT;", plan, name)
}

// AddMigration adds a migration to the end of the plan file
func (plan *PlanFile) AddMigration(config *Config, name string, comment string) error {
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
