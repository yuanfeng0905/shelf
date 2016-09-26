// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/cayleygraph/cayley"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/coralproject/shelf/internal/xenia"
	"github.com/coralproject/shelf/internal/xenia/query"
)

// execHandle maintains the set of handlers for the exec api.
type execHandle struct{}

// Exec fronts the access to the exec service functionality.
var Exec execHandle

//==============================================================================

// Name runs the specified Set and return results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) Name(c *app.Context) error {
	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	return execute(c, set)
}

// NameOnView runs the specified Set on a view and return results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) NameOnView(c *app.Context) error {

	// Retrieve the query set.
	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	// Execute the view.
	viewParams := wire.ViewParams{
		ViewName:          c.Params["view"],
		ItemKey:           c.Params["item"],
		ResultsCollection: set.Collection,
	}

	if _, err := wire.Execute(c.SessionID, c.Ctx["DB"].(*db.DB), c.Ctx["Graph"].(*cayley.Handle), viewParams); err != nil {
		return err
	}

	// Execute the query.
	if err := execute(c, set); err != nil {
		return err
	}

	return nil
}

// Custom runs the provided Set and return results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) Custom(c *app.Context) error {
	var set *query.Set
	if err := json.NewDecoder(c.Request.Body).Decode(&set); err != nil {
		return err
	}

	return execute(c, set)
}

// CustomOnView runs the provided Set on a view and return results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) CustomOnView(c *app.Context) error {

	// Decode the query set.
	var set *query.Set
	if err := json.NewDecoder(c.Request.Body).Decode(&set); err != nil {
		return err
	}

	// Execute the view.
	viewParams := wire.ViewParams{
		ViewName:          c.Params["view"],
		ItemKey:           c.Params["item"],
		ResultsCollection: set.Collection,
	}

	if _, err := wire.Execute(c.SessionID, c.Ctx["DB"].(*db.DB), c.Ctx["Graph"].(*cayley.Handle), viewParams); err != nil {
		return err
	}

	// Execute the query.
	if err := execute(c, set); err != nil {
		return err
	}

	return nil
}

//==============================================================================

// execute takes a context and Set and executes the set returning
// any possible response.
func execute(c *app.Context, set *query.Set) error {
	var vars map[string]string
	if c.Request.URL.RawQuery != "" {
		if m, err := url.ParseQuery(c.Request.URL.RawQuery); err == nil {
			vars = make(map[string]string)
			for k, v := range m {
				vars[k] = v[0]
			}
		}
	}

	result := xenia.Exec(c.SessionID, c.Ctx["DB"].(*db.DB), set, vars)

	c.Respond(result, http.StatusOK)
	return nil
}
