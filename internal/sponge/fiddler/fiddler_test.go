package fiddler_test

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/sponge/fiddler"
	"github.com/coralproject/shelf/internal/sponge/fiddler/fiddlerfix"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("SHELF")
}

// prefix is what we are looking to delete after the test.
const prefix = "ITEST_"

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
