// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/xenia/query"
)

// queryHandle maintains the set of handlers for the query api.
type queryHandle struct{}

// Query fronts the access to the query service functionality.
var Query queryHandle

//==============================================================================

// List returns all the existing Set names in the system.
// 200 Success, 404 Not Found, 500 Internal
func (queryHandle) List(c *app.Context) error {
	sets, err := query.GetAll(c.SessionID, c.Ctx["DB"].(*db.DB), nil)
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(sets, http.StatusOK)
	return nil
}

// Retrieve returns the specified Set from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (queryHandle) Retrieve(c *app.Context) error {
	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(set, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted Set document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (queryHandle) Upsert(c *app.Context) error {
	var set query.Set
	if err := json.NewDecoder(c.Request.Body).Decode(&set); err != nil {
		return err
	}

	if err := query.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &set); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

// EnsureIndexes makes sure indexes for the specified set exist.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (queryHandle) EnsureIndexes(c *app.Context) error {
	db := c.Ctx["DB"].(*db.DB)

	set, err := query.GetByName(c.SessionID, db, c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	if err := query.EnsureIndexes(c.SessionID, db, set); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified Set from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (queryHandle) Delete(c *app.Context) error {
	if err := query.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"]); err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
