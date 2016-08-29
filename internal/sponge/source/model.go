package source

import (
	"database/sql"

	"github.com/coralproject/pillar/pkg/db"
	validator "gopkg.in/bluesuncorp/validator.v8"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Sourcer interface to get data from external sources
type Sourcer interface {
	GetData(string, *Options) (db.Iterator, error)
}

//==============================================================================

type MSQL struct {
	DB *sql.DB
}

//==============================================================================

type PSQL struct {
	DB *sql.DB
}

//==============================================================================
