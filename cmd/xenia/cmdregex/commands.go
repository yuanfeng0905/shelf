package cmdregex

import (
	"github.com/ardanlabs/kit/db"
	"github.com/spf13/cobra"
)

// regexCmd represents the parent for all regex cli commands.
var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "regex provides a xenia CLI for managing regexs.",
}

// conn holds the session for the DB access.
var conn *db.DB

// GetCommands returns the regex commands.
func GetCommands(db *db.DB) *cobra.Command {
	conn = db
	addUpsert()
	addGet()
	addDel()
	addList()
	return regexCmd
}
