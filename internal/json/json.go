package json

import (
	"encoding/json"
	"os"
)

type NoteStruct struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

const notesFile = "./data/notes.json"

func ReadNotes() ([]NoteStruct, error) {
	file, err := os.Open(notesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []NoteStruct{}, nil // Return an empty list if the file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	var notes []NoteStruct
	err = json.NewDecoder(file).Decode(&notes)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func WriteNotes(notes []NoteStruct) error {
	file, err := os.Create(notesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	return encoder.Encode(notes)
}
