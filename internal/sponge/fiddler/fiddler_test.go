package fiddler_test

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/sponge/fiddler"
	"github.com/coralproject/shelf/internal/sponge/fiddler/fiddlerfix"
	"github.com/coralproject/shelf/internal/sponge/item"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("SHELF")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	cfg := mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(cfg)
}

// prefix is what we are looking to delete after the test.
const prefix = "FTEST_"

// TestTransform tests if we can transform a row into an item.
func TestTransform(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	strategyFile := os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/fiddler/fiddlerfix/strategy_coral_test.json"

	err := fiddler.Setup(tests.Context, strategyFile)
	if err != nil {
		t.Fatalf("\t%s\tShould be able retrieve fiddler fixture : %s", tests.Failed, err)
	}

	t.Log("Given the need to transform data into an item.")
	{
		t.Log("\tWhen starting from an existed set of data")
		{
			//----------------------------------------------------------------------
			// Get the fixture.
			row, err := fiddlerfix.GetRawDataRow()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve fiddler fixture : %s", tests.Failed, err)
			}

			entityName := "comments"

			//----------------------------------------------------------------------
			// Transform into an item.

			item, err := fiddler.Transform(tests.Context, row, entityName)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to transform into an item : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to transform into an item.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the item we expected.

			if item.ID != row["commentID"] || item.Type != "comments" || item.Version != 1 {
				t.Logf("\t%+v", item.ID)
				t.Logf("\t%+v", item.Type)
				t.Logf("\t%+v", item.Version)
				t.Fatalf("\t%s\tShould be able to get back the right item.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the right item.", tests.Success)

			if item.Data.(map[string]interface{})["commentTitle"] != "Titulo" {
				t.Logf("\t%+v", item.Data)
				t.Fatalf("\t%s\tShould be able to store the raw data.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to store the raw data.", tests.Success)
		}
	}
}

// TestBulkTransform tests if we can transform and insert items
func TestBulkTransform(t *testing.T) {

	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := fiddlerfix.Remove(tests.Context, db, prefix); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the items : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the items.", tests.Success)

		if err := fiddlerfix.RemoveFixtures(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the fixtures : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the fixtures.", tests.Success)
	}()

	strategyFile := os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/fiddler/fiddlerfix/strategy_coral_test.json"

	err = fiddler.Setup(tests.Context, strategyFile)
	if err != nil {
		t.Fatalf("\t%s\tShould be able retrieve fiddler fixture : %s", tests.Failed, err)
	}

	t.Log("Given the need to import transformed data into the system.")
	{
		t.Log("\tWhen starting from an existed set of data")
		{
			//----------------------------------------------------------------------
			// Get the fixture.
			iter, err := fiddlerfix.GetRawDataIterator(tests.Context, db, prefix)
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve fiddler fixture : %s", tests.Failed, err)
			}

			entityName := "comments"

			//----------------------------------------------------------------------
			// Transform and Import.

			err = fiddler.BulkTransform(tests.Context, db, iter, entityName)

			if err != nil {
				t.Fatalf("\t%s\tShould be able to transform into an item : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to transform into an item.", tests.Success)

			//----------------------------------------------------------------------
			// Get the items.

			itemsBack, err := item.GetByIDs(tests.Context, db, []string{"16546089", "16546090"})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the item by ID : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the item by ID.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the item we expected.

			expectedItem := item.Item{
				ID:      "16546089",
				Type:    "comments",
				Version: 1,
				Data: map[string]string{
					"commentID":    "16546089",
					"assetID":      "3441001",
					"statusID":     "3",
					"commentTitle": "<br/>",
				},
			}

			if expectedItem.ID != itemsBack[0].ID || expectedItem.Type != itemsBack[0].Type || expectedItem.Version != itemsBack[0].Version {
				t.Logf("\t%+v", itemsBack[0])
				t.Fatalf("\t%s\tShould be able to get back the same item.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same item.", tests.Success)

		}
	}
}
