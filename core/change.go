package core

import (
	"regexp"
	"strings"
	"time"
)

type Change struct {
	Name      string
	CreatedBy string
	Date      time.Time
	Comment   string
	Content   string
	Project   string
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
var ChangeRegex = regexp.MustCompile(`(?P<name>[^\s]+)\s+(?P<date>[^\s]+)\s+(?P<created_by>[^#]+)#?(?P<comment>.*)`)

func ParseChange(raw string) (Change, error) {
	if match := ChangeRegex.FindStringSubmatch(raw); match != nil {
		date, err := time.Parse(time.RFC3339, strings.TrimSpace(match[2]))

		if err != nil {
			return Change{}, err
		}

		return Change{
			Name:      strings.TrimSpace(match[1]),
			Date:      date,
			CreatedBy: strings.TrimSpace(match[3]),
			Comment:   strings.TrimSpace(match[4]),
			Content:   "",
		}, nil
	}

	return Change{}, &ParseError{
		Value:   raw,
		Message: "Value does not match Change format.",
	}
}
