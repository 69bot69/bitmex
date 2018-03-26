package writer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"bitmex.lib/models"
)

func WriteJson(a []models.Indicator) error {
	j, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}

	t := time.Now()
	s := t.String()

	if err := ioutil.WriteFile(fmt.Sprintf("./data/indicator_%s.json", s), j, 0644); err != nil {
		return err
	}

	return nil
}
