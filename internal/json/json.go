package json

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type NoteStruct struct {
	ID      int64  `json:"id"` // TODO
	Title   string `json:"title"`
	Content string `json:"content"`
}

const notesFile = "data/notes.json"

// ReadNotes reads the notes from the JSON file.
func ReadNotes() ([]NoteStruct, error) {
	file, err := os.Open(notesFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Return an empty slice if the file doesn't exist
			return []NoteStruct{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var notes []NoteStruct
	if err = json.NewDecoder(file).Decode(&notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// WriteNotes writes the notes to the JSON file.
func WriteNotes(notes []NoteStruct) error {
	// Ensure the directory exists
	dir := filepath.Dir(notesFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(notesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	if err = encoder.Encode(notes); err != nil {
		return err
	}

	return nil
}
