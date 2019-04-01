package core

import (
	"regexp"
	"strings"
	"time"
)

type Migration struct {
	Name      string
	CreatedBy string
	Date      time.Time
	Comment   string
	Content   string
}

// ParseError describes a problem parsing a time string.
type ParseError struct {
	Value   string
	Message string
}

// Error returns the string representation of a ParseError.
func (e *ParseError) Error() string {
	return e.Message
}

//clinical/schema/patient 2019-04-01T00:47:14Z Joe Andaverde <joe.andaverde@gmail.com> # Note
var migrationRegex = regexp.MustCompile(`(?P<name>[^\s]+)\s+(?P<date>[^\s]+)\s+(?P<created_by>[^#]+)#?(?P<comment>.*)`)

func ParseMigration(raw string) (Migration, error) {
	if match := migrationRegex.FindStringSubmatch(raw); match != nil {
		date, err := time.Parse(time.RFC3339, strings.TrimSpace(match[2]))

		if err != nil {
			return Migration{}, err
		}

		return Migration{
			Name:      strings.TrimSpace(match[1]),
			Date:      date,
			CreatedBy: strings.TrimSpace(match[3]),
			Comment:   strings.TrimSpace(match[4]),
			Content:   "",
		}, nil
	}

	return Migration{}, &ParseError{
		Value:   raw,
		Message: "Value does not match migration format.",
	}
}
