package regex_test

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/coralproject/xenia/pkg/regex"
	"github.com/coralproject/xenia/pkg/regex/rfix"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("XENIA")

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

//==============================================================================

// TestUpsertCreateRegex tests if we can create a regex record in the db.
func TestUpsertCreateRegex(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := rfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the regex : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the regex.", tests.Success)
	}()

	t.Log("Given the need to save a regex into the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			if _, err = regex.GetLastHistoryByName(tests.Context, db, rgx1.Name); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the regex from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the regex from history.", tests.Success)

			rgx2, err := regex.GetByName(tests.Context, db, rgx1.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the regex.", tests.Success)

			if rgx1.Compile, err = regexp.Compile(rgx1.Expr); err != nil {
				t.Fatalf("\t%s\tShould be able to compile the regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to compile the regex.", tests.Success)

			if !reflect.DeepEqual(*rgx1, *rgx2) {
				t.Logf("\t%+v", rgx1)
				t.Logf("\t%+v", rgx2)
				t.Errorf("\t%s\tShould be able to get back the same regex values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same regex values.", tests.Success)
			}
		}
	}
}

// TestGetRegexNames validates retrieval of Regex record names.
func TestGetRegexNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	rgxName := "RTEST_basic"

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := rfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the regexs : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the regexs.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of regexs.")
	{
		t.Log("\tWhen using two regexs")
		{
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			rgx2 := *rgx1
			rgx2.Name += "2"
			if err := regex.Upsert(tests.Context, db, &rgx2); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second regex.", tests.Success)

			names, err := regex.GetNames(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the regex names : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the regex names", tests.Success)

			var count int
			for _, name := range names {
				if len(name) > 5 && name[0:5] == "RTEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two regexs : %d", tests.Failed, len(names))
			}
			t.Logf("\t%s\tShould have two regexs.", tests.Success)

			if !strings.Contains(names[0], rgxName) || !strings.Contains(names[1], rgxName) {
				t.Errorf("\t%s\tShould have \"%s\" in the name : %s", tests.Failed, rgxName, names[0])
			} else {
				t.Logf("\t%s\tShould have \"%s\" in the name.", tests.Success, rgxName)
			}
		}
	}
}

// TestGetRegexs validates retrieval of all Regex records.
func TestGetRegexs(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := rfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the regexs : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the regexs.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of regexs.")
	{
		t.Log("\tWhen using two regexs")
		{
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			rgx1.Name += "2"
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second regex.", tests.Success)

			rgxs, err := regex.GetRegexs(tests.Context, db, nil)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the regexs : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the regexs", tests.Success)

			var count int
			for _, rgx := range rgxs {
				if len(rgx.Name) > 5 && rgx.Name[0:5] == "RTEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two regexs : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two regexs.", tests.Success)
		}
	}
}

// TestGetRegexByNames validates retrieval of Regex records by a set of names.
func TestGetRegexByNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := rfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the regexs : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the regexs.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of regex values.")
	{
		t.Log("\tWhen using two regexs")
		{
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			rgx2 := *rgx1
			rgx2.Name += "2"
			if err := regex.Upsert(tests.Context, db, &rgx2); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second regex.", tests.Success)

			regexs, err := regex.GetByNames(tests.Context, db, []string{rgx1.Name, rgx2.Name})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the regexs by names : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the regexs by names", tests.Success)

			var count int
			for _, rgx := range regexs {
				if len(rgx.Name) > 5 && rgx.Name[0:5] == "RTEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two regexs : %d", tests.Failed, len(regexs))
			}
			t.Logf("\t%s\tShould have two regexs.", tests.Success)

			if regexs[0].Name != rgx1.Name || regexs[1].Name != rgx2.Name {
				t.Errorf("\t%s\tShould have retrieve the correct regexs.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have retrieve the correct regexs.", tests.Success)
			}
		}
	}
}

// TestGetLastRegexHistoryByName validates retrieval of Regex from the history
// collection.
func TestGetLastRegexHistoryByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	rgxName := "RTEST_basic"

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := rfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the regexs : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the regexs.", tests.Success)
	}()

	t.Log("Given the need to retrieve a regex from history.")
	{
		t.Log("\tWhen using regex", rgx1)
		{
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			rgx1.Expr = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"

			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			rgx2, err := regex.GetLastHistoryByName(tests.Context, db, rgxName)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the last regex from history : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the last regex from history.", tests.Success)

			if !reflect.DeepEqual(*rgx1, *rgx2) {
				t.Logf("\t%+v", rgx1)
				t.Logf("\t%+v", rgx2)
				t.Errorf("\t%s\tShould be able to get back the same regex values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same regex values.", tests.Success)
			}
		}
	}
}

// TestUpsertUpdateRegex validates update operation of a given Regex.
func TestUpsertUpdateRegex(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := rfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the regexs : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the regexs.", tests.Success)
	}()

	t.Log("Given the need to update a regex into the database.")
	{
		t.Log("\tWhen using two regexs")
		{
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			rgx2 := *rgx1
			rgx2.Expr = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"

			if err := regex.Upsert(tests.Context, db, &rgx2); err != nil {
				t.Fatalf("\t%s\tShould be able to update a regex record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a regex record.", tests.Success)

			if _, err := regex.GetLastHistoryByName(tests.Context, db, rgx1.Name); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the regex from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the regex from history.", tests.Success)

			updRgx, err := regex.GetByName(tests.Context, db, rgx2.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a regex record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a regex record.", tests.Success)

			if updRgx.Name != rgx1.Name {
				t.Errorf("\t%s\tShould be able to get back the same regex name.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same regex name.", tests.Success)
			}

			if updRgx.Expr == rgx1.Expr {
				t.Logf("\t%+v", updRgx.Expr)
				t.Logf("\t%+v", rgx1.Expr)
				t.Errorf("\t%s\tShould have an updated regex record.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have an updated regex record.", tests.Success)
			}
		}
	}
}

// TestDeleteRegex validates the removal of a regex from the database.
func TestDeleteRegex(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	rgxName := "RTEST_basic"
	rgxBadName := "RTEST_basic_advice"

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := rfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the regexs : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the regexs.", tests.Success)
	}()

	t.Log("Given the need to delete a regex in the database.")
	{
		t.Log("\tWhen using regex", rgx1)
		{
			if err := regex.Upsert(tests.Context, db, rgx1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a regex : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a regex.", tests.Success)

			if err := regex.Delete(tests.Context, db, rgxName); err != nil {
				t.Fatalf("\t%s\tShould be able to delete a regex using its name[%s]: %s", tests.Failed, rgxName, err)
			}
			t.Logf("\t%s\tShould be able to delete a regex using its name[%s]:", tests.Success, rgxName)

			if err := regex.Delete(tests.Context, db, rgxBadName); err == nil {
				t.Fatalf("\t%s\tShould not be able to delete a regex using wrong name name[%s]", tests.Failed, rgxBadName)
			}
			t.Logf("\t%s\tShould not be able to delete a regex using wrong name name[%s]", tests.Success, rgxBadName)

			if _, err := regex.GetByName(tests.Context, db, rgxName); err == nil {
				t.Fatalf("\t%s\tShould be able to validate regex with Name[%s] does not exists: %s", tests.Failed, rgxName, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate regex with Name[%s] does not exists:", tests.Success, rgxName)
		}
	}
}

// TestAPIFailureRegexs validates the failure of the api using a nil session.
func TestAPIFailureRegexs(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	rgxName := "RTEST_unknown"

	const fixture = "basic.json"
	rgx1, err := rfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load regex record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load regex record from file.", tests.Success)

	t.Log("Given the need to validate failure of API with bad session.")
	{
		t.Log("When giving a nil session")
		{
			err := regex.Upsert(tests.Context, nil, rgx1)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused create by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused create by api with bad session: %s", tests.Success, err)

			_, err = regex.GetNames(tests.Context, nil)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = regex.GetByName(tests.Context, nil, rgxName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = regex.GetByNames(tests.Context, nil, nil)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = regex.GetLastHistoryByName(tests.Context, nil, rgxName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			err = regex.Delete(tests.Context, nil, rgxName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused delete by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused delete by api with bad session: %s", tests.Success, err)
		}
	}
}
