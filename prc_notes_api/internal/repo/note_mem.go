package repo

import (
	"fmt"
	"sync"
	"time"

	"example.com/prc_notes_api/internal/core"
)

type NoteRepoMem struct {
	mu    sync.Mutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}

func (r *NoteRepoMem) Create(n core.ReqNote) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	now := time.Now()
	note := &core.Note{
		Title:     n.Title,
		Content:   n.Content,
		ID:        r.next,
		CreatedAt: now,
		UpdatedAt: &now,
	}
	r.notes[note.ID] = note
	return note, nil
}

func (r *NoteRepoMem) Get(id int64) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, exists := r.notes[id]

	if !exists {
		return nil, fmt.Errorf("note with id %d does not exist", id)
	}

	return n, nil
}

func (r *NoteRepoMem) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.notes, id)
	return nil
}

func (r *NoteRepoMem) Update(id int64, updNote *core.ReqNote) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.notes[id].Content = updNote.Content
	r.notes[id].UpdatedAt = &now
	r.notes[id].Title = updNote.Title

	return nil
}

func (r *NoteRepoMem) Notes() map[int64]*core.Note {
	return r.notes
}
