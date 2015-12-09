package query_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/coralproject/shelf/pkg/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	tests.InitMGO()
}

//==============================================================================

// TestUmarshalMongoScript tests the ability to convert string based Mongo
// commands into a bson map for processing.
func TestUmarshalMongoScript(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	scripts := []struct {
		text string
		qry  *query.Query
		cmp  bson.M
	}{
		{
			`{"name":"bill"}`,
			nil,
			bson.M{"name": "bill"},
		},
		{
			`{"date":"ISODate('2013-01-16T00:00:00.000Z')"}`,
			&query.Query{HasDate: true},
			bson.M{"date": time.Date(2013, time.January, 16, 0, 0, 0, 0, time.UTC)},
		},
		{
			`{"_id":"ObjectId(\"5660bc6e16908cae692e0593\")"}`,
			&query.Query{HasObjectID: true},
			bson.M{"_id": bson.ObjectIdHex("5660bc6e16908cae692e0593")},
		},
	}

	t.Logf("Given the need to convert mongo commands.")
	{
		for _, script := range scripts {
			t.Logf("\tWhen using %s with %+v", script.text, script.qry)
			{
				b, err := query.UmarshalMongoScript(script.text, script.qry)
				if err != nil {
					t.Errorf("\t%s\tShould be able to convert without an error : %v", tests.Failed, err)
					continue
				}
				t.Logf("\t%s\tShould be able to convert without an error.", tests.Success)

				if eq := compareBson(b, script.cmp); !eq {
					t.Log(b)
					t.Log(script.cmp)
					t.Errorf("\t%s\tShould get back the expected bson document.", tests.Failed)
					continue
				}
				t.Logf("\t%s\tShould get back the expected bson document.", tests.Success)
			}
		}
	}
}

// TestExecuteSet tests the execution of different Sets that should succeed.
func TestExecuteSet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	execSet := getExecSet()

	db := db.NewMGO()
	defer db.CloseMGO()

	t.Logf("Given the need to execute mongo commands.")
	{
		err := generateTestData(db)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to load system with test data : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to load system with test data.", tests.Success)

		defer dropTestData()

		for i, es := range execSet {
			t.Logf("\tWhen using Execute Set %d", i)
			{
				result := query.ExecuteSet(tests.Context, db, es.set, es.vars)
				if result.Error {
					t.Errorf("\t%s\tShould be able to execute the query set : %+v", tests.Failed, result.Results)
					continue
				}
				t.Logf("\t%s\tShould be able to execute the query set.", tests.Success)

				data, err := json.Marshal(result)
				if err != nil {
					t.Errorf("\t%s\tShould be able to marshal the result : %s", tests.Failed, err)
					continue
				}
				t.Logf("\t%s\tShould be able to marshal the result.", tests.Success)

				var res query.Result
				if err := json.Unmarshal(data, &res); err != nil {
					t.Errorf("\t%s\tShould be able to unmarshal the result : %s", tests.Failed, err)
					continue
				}
				t.Logf("\t%s\tShould be able to unmarshal the result.", tests.Success)

				if string(data) != es.result {
					t.Errorf("\t%s\tShould have the correct result.", tests.Failed)
					continue
				}
				t.Logf("\t%s\tShould have the correct result", tests.Success)
			}
		}
	}
}

//==============================================================================

// execSet represents the table for the table test of execution tests.
type execSet struct {
	set    *query.Set
	vars   map[string]string
	result string
}

// docs represents what a user will receive after
// excuting a successful set.
type docs struct {
	Name string
	Docs []bson.M
}

// getExecSet returns the table for the testing.
func getExecSet() []execSet {
	sets := make([]execSet, 1)
	sets[0].set, sets[0].result = querySetBasic()

	return sets
}

// querySetBasic starts with a simple query set.
func querySetBasic() (*query.Set, string) {
	set := query.Set{
		Name:    "test",
		Enabled: true,
		Queries: []query.Query{
			{
				Name:       "Q1",
				Type:       "pipeline",
				Collection: "test_query",
				Return:     true,
				Scripts: []string{
					`{"$match" : {"station_id" : "42021"}}`,
					`{"$project": {"_id": 0, "name": 1}}`,
				},
			},
		},
	}

	result := `{"results":[{"Name":"Q1","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`

	return &set, result
}

//==============================================================================

// generateTestData creates a temp collection with data
// that can be used for testing things.
func generateTestData(db *db.DB) error {
	file, err := os.Open("exec_test_data.json")
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var rawDocs []bson.M
	if err := json.Unmarshal(data, &rawDocs); err != nil {
		return err
	}

	var docs []interface{}
	for _, rd := range rawDocs {
		mar, err := json.Marshal(rd)
		if err != nil {
			return err
		}

		doc, err := query.UmarshalMongoScript(string(mar), &query.Query{HasDate: true})
		if err != nil {
			return err
		}

		docs = append(docs, doc)
	}

	f := func(c *mgo.Collection) error {
		return c.Insert(docs...)
	}

	if err := db.ExecuteMGO(tests.Context, "test_query", f); err != nil {
		return err
	}

	return nil
}

// dropTestData drops the temp collection.
func dropTestData() {
	db := db.NewMGO()
	defer db.CloseMGO()

	mongo.GetCollection(db.MGOConn, "test_query").DropCollection()
}

// compareBson compares two bson maps for equivalence.
func compareBson(m1 bson.M, m2 bson.M) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}

	for k, v := range m2 {
		if m1[k] != v {
			return false
		}
	}

	return true
}
