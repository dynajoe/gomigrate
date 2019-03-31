package cmds

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var AddMigrationDoc = `
Perform a GET request to the value of the --url flag. If no url is specified
then http://www.google.com is used as the url.
`

var urlFlag string

func AddMigration() error {
	res, err := http.Get(urlFlag)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

func AddMigrationFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&urlFlag, "url", "http://www.google.com", "Page to download")
}
