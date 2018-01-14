package handlers

import (
	"github.com/h3poteto/fascia/server/commands/board"
)

// ListOptionAll returns all list options
func ListOptionAll() ([]*board.ListOption, error) {
	return board.ListOptionAll()
}
