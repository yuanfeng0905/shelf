// Package fiddler transform, through a strategy file, data from external source into our local coral schema.
//
// It gets a map[string]interface{} as a row and the coral's model that is going to convert it to.
// With that model goes to the strategy file to see what is the transformation that needs to do.
package fiddler

import (
	"fmt"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/sponge/strategy"
)

// global variables related to strategy
var (
	s *strategy.Strategy
)

// Setup sets the strategy that we are going to use for the fiddler's transformation
func Setup(context interface{}, strategyFile string) error {
	var err error

	if s, err = strategy.New(context, strategyFile); err != nil {
		log.Error(context, "Started", err, "Fail on starting strategy")
		return err
	}

	return nil
}

// =============================================================================

//func BulkTransform(context interface{},  *mgo.Iter) (*mgo.Bulk, error) {}

// Transform transforms a row of data into the coral schema
func Transform(context interface{}, row map[string]interface{}, entityName string) (*item.Item, error) { //transformation, error

	var err error

	entity, err := s.GetEntity(context, entityName)
	if err != nil {
		return nil, err
	}

	if entity.Local == "" {
		errLocalNotFound := fmt.Errorf("No local value for entity %v found in the strategy file.", entityName)
		log.Error(context, "Transform", errLocalNotFound, "Not Found.")
		return nil, errLocalNotFound
	}

	idValue, ok := row[entity.IDField]
	if !ok {
		errIDFieldNotFound := fmt.Errorf("No local value for ID field %v found in the data.", entity.IDField)
		log.Error(context, "Transform", errIDFieldNotFound, "Not Found.")
		return nil, errIDFieldNotFound
	}

	i := new(item.Item)

	i.ID = idValue.(string)
	i.Type = entity.Local
	i.Version = s.Version
	i.Data = row

	return i, nil
}
