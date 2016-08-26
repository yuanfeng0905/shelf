package fiddlerfix

import (
	"encoding/json"
	"os"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/fiddler/fiddlerfix/"
}

// GetRawDataRow returns raw data.
func GetRawDataRow() (map[string]interface{}, error) {
	file, err := os.Open(path + "comments_rawdata.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
