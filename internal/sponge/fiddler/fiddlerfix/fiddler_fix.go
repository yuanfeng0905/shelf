package fiddlerfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/sponge/item"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/fiddler/fiddlerfix/"
}

// GetRawDataRow returns raw data.
func GetRawDataRow() (map[string]interface{}, error) {
	file, err := os.Open(path + "comment_rawdata.json")
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

// GetRawData returns raw data.
func GetRawDataIterator(context interface{}, db *db.DB, prefix string) (*mgo.Iter, error) { // TO DO: RETURN DB.ITER , ITERATOR FOR ANY DBMS
	file, err := os.Open(path + "comments_rawdata.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	// Insert the fixtures into a temporal db's collection. Can we do this with db.Iter?
	var b *mgo.Bulk
	b, err = db.BulkOperationMGO(context, "rawdata")
	if err != nil {
		return nil, err
	}

	b.Unordered()

	// insert data into temporal Collection
	b.Insert(data...)

	_, err = b.Run()
	if err != nil {
		return nil, err
	}

	q := bson.M{}

	iter, err := db.BatchedQueryMGO(context, "rawdata", q)
	if err != nil {
		return nil, err
	}

	return iter, nil
}

func Remove(context interface{}, db *db.DB, prefix string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"item_id": bson.RegEx{Pattern: prefix}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, item.Collection, f); err != nil {
		return err
	}

	return nil
}

func RemoveFixtures(context interface{}, db *db.DB) error {
	f := func(c *mgo.Collection) error {
		_, err := c.RemoveAll(bson.M{})
		return err
	}

	if err := db.ExecuteMGO(context, "rawdata", f); err != nil {
		return err
	}

	return nil
}
