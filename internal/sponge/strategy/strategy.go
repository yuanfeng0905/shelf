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

	"github.com/ardanlabs/kit/log"
	"github.com/xeipuuv/gojsonschema"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Strategy Not found")
	ErrNotValid = errors.New("Strategy Not valid")
)

// New creates a new strategy struct variable from the json file
func New(context interface{}, strategyFile string) (*Strategy, error) {
	log.Dev(context, "New", "Started: ", strategyFile)

	// validate Strategy file
	if ok, err := Validate(context, strategyFile); !ok {
		return nil, err
	}

	strategy, err := Read(context, strategyFile)
	if err != nil {
		return nil, err
	}

	return strategy, nil
}

// IsEmpty check if the string is empty
func IsEmpty(context interface{}, fileName string) error {
	if _, err := os.Stat(fileName); err != nil {
		log.Error(context, "IsEmpty", err, "Completed")
		return err
	}

	return nil
}

// Validate the strategy file that is loaded into STRATEGY_CONF environment variable.
func Validate(context interface{}, strategyFile string) (bool, error) {

	schemaFile := "file:///" + os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/strategy/schema_strategy.json"
	if err := IsEmpty(context, strategyFile); err != nil {
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
		log.Error(context, "Validate", err, "Completed")
		return false, errNotValid
	}

	if !result.Valid() {
		verror := fmt.Sprintf("Validation Errors on %s:\n", strategyFile)
		for _, err := range result.Errors() {
			verror = verror + fmt.Sprintf("%v - %s \n", err.Details(), err.Description())
		}
		err = errors.New(verror)
		log.Error(context, "Validate", err, "Completed")
		return false, err
	}

	return true, nil
}

// Read the strategy file and do the validation into the Strategy struct
func Read(context interface{}, f string) (*Strategy, error) {

	var strategy *Strategy

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Error(context, "Validate", err, "Completed")
		return nil, err
	}

	err = json.Unmarshal(content, &strategy)

	return strategy, err
}
