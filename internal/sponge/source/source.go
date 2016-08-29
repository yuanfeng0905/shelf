/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* PostgreSQL
* MongoDB
* Webservice

*/

package source

import (
	"errors"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/sponge/strategy"
)

type Options struct {
	Limit                 int
	Offset                int
	Orderby               string
	Query                 string
	Types                 string
	Importonlyfailed      bool
	ReportOnFailedRecords bool
	Reportdbfile          string
	TimeWaiting           int
}

var (
	s strategy.Strategy
)

var ErrNotFound = errors.New("Source Not Found.")

func New(source string) (Sourcer, error) {

	switch source {
	// case "mysql":
	// 	return
	case "mongodb":
		// Get MongoDB connection string
		return *db.DB, nil
	}
	return nil, ErrNotFound
}

// GetForeignEntity returns the name of the foreign source's entity
func GetForeignEntity(context interface{}, local string) (string, err) {

	e, err := s.GetEntity(context, local)
	if err != nil {
		log.Error(context, "GetForeignEntity", err)
		return "", err
	}
	return e.Foreign, nil
}
