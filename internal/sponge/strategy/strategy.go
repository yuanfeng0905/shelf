/*
Package strategy handles the loading and distribution of configuration related with mapping external data to our own schema.
*/
package strategy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Strategy Not found")
	ErrNotValid = errors.New("Strategy Not valid")
)

// Field holds the specific transformation. Foreign converts into a Local through the specific transformation. Type is the data type
type Field struct {
	Foreign        string `bson:"foreign" json:"foreign"`
	Local          string `bson:"local" json:"local"`
	Transformation string `bson:"transformation" json:"transformation"`
	Type           string `bson:"type" json:"type"`
	DateTimeFormat string `jbson:"datetimeformat" json:"datetimeformat"`
}

// Entity holds the struct on what is the external source's entity name and fields
type Entity struct {
	Foreign        string            `json:"foreign"`
	Local          string            `json:"local"`
	OrderBy        string            `json:"orderby"`
	Fields         []Field           `json:"fields"`
	DateTimeFormat string            `json:"datetimeformat"`
	Status         map[string]string `json:"status"`
}

// Strategy explains which entities or data we are getting from the source and which transformation nees to happen.
type Strategy struct {
	Name           string            `bson:"name" json:"name"`
	DateTimeFormat string            `bson:"datetimeformat" json:"datetimeformat"`
	Entities       map[string]Entity `bson:"entities" json:"entities"`
}

// =============================================================================

// New creates a new strategy struct variable from the json file
func New() (*Strategy, error) {

	//read STRATEGY_CONF env variable
	strategyFile := os.Getenv("STRATEGY_CONF")

	// validate Strategy file
	if ok, err := Validate(strategyFile); !ok {
		return nil, err
	}

	strategy, err := Read(strategyFile)
	if err != nil {
		return nil, err
	}

	return strategy, nil
}

// IsEmpty check if the string is empty
func IsEmpty(fileName string) error {
	_, err := os.Stat(fileName)
	// log
	return err
}

// Validate the strategy file that is loaded into STRATEGY_CONF environment variable.
func Validate(strategyFile string) (bool, error) {

	schemaFile := "file:///" + os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/strategy/schema_strategy.json"
	if err := IsEmpty(strategyFile); err != nil {
		return false, err
	}

	schemaLoader := gojsonschema.NewReferenceLoader(schemaFile)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + strategyFile)

	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return false, fmt.Errorf("Not able to load the schema for %s. Error: %v", schemaFile, err)
	}

	result, err := schema.Validate(documentLoader)
	if err != nil {
		errNotValid := fmt.Errorf("%s strategy not Valid: %v", strategyFile, err)
		//log
		return false, errNotValid
	}

	if !result.Valid() {
		verror := fmt.Sprintf("Validation Errors on %s:\n", strategyFile)
		for _, err := range result.Errors() {
			verror = verror + fmt.Sprintf("%v - %s \n", err.Details(), err.Description())
		}
		//log
		return false, errors.New(verror)
	}

	return true, nil
}

// Read the strategy file and do the validation into the Strategy struct
func Read(f string) (*Strategy, error) {

	var strategy *Strategy

	content, err := ioutil.ReadFile(f)
	if err != nil {
		// log
		return nil, err
	}

	err = json.Unmarshal(content, &strategy)

	return strategy, err
}

// =============================================================================

// HasArrayField returns true if the entity has fields that are type array and need to be loop through
func (s Strategy) HasArrayField(e Entity) bool {

	for _, f := range e.Fields {
		if f.Type == "Array" {
			return true
		}
	}
	return false
}

// GetDefaultDateTimeFormat gets the default datetime format
func (s Strategy) GetDefaultDateTimeFormat() string {
	return s.DateTimeFormat
}

// GetDateTimeFormat returns the datetime format for this strategy
func (s Strategy) GetDateTimeFormat(entity string, field string) string {

	for _, f := range s.Entities[entity].Fields {
		if f.Local == field {
			return f.DateTimeFormat
		}
	}
	return s.GetDefaultDateTimeFormat()
}

// GetEntities returns a list of the entities defined in the transformations file
func (s Strategy) GetEntities() map[string]Entity {
	return s.Entities
}

// GetEntityForeignName returns the external source's entity mapped to the coral model
func (s Strategy) GetEntityForeignName(coralName string) string {
	return s.Entities[coralName].Foreign
}

// GetEntityForeignFields returns the external source's entity fields mapped to the coral model
func (s Strategy) GetEntityForeignFields(coralName string) []Field {
	return s.Entities[coralName].Fields
}

// GetOrderBy returns the order by field definied in the transformations file
func (s Strategy) GetOrderBy(coralName string) string {
	return s.Entities[coralName].OrderBy
}

// GetStatus returns the mapping of the external status into the coral one
func (s Strategy) GetStatus(coralName string, foreign string) string {
	return s.Entities[coralName].Status[foreign]
}
