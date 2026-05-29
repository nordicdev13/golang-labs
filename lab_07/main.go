package main

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type InMemoryStore struct {
	mu     sync.RWMutex
	notes  []Note
	nextID int
}

func main() {
	store := &InMemoryStore{
		notes: []Note{
			{ID: 1, Title: "Buy milk", Content: "2 liters"},
		},
		nextID: 2,
	}

	app := fiber.New()

	app.Get("/notes", store.GetAllNotes)
	app.Get("/notes/:id", store.GetNoteByID)
	app.Post("/notes", store.CreateNote)
	app.Put("/notes/:id", store.UpdateNote)
	app.Delete("/notes/:id", store.DeleteNote)

	app.Listen(":8080")
}

func (s *InMemoryStore) GetAllNotes(c *fiber.Ctx) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.notes) == 0 {
		return c.JSON([]Note{})
	}

	return c.JSON(s.notes)
}

func (s *InMemoryStore) GetNoteByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, note := range s.notes {
		if note.ID == id {
			return c.JSON(note)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note not found"})
}

func (s *InMemoryStore) CreateNote(c *fiber.Ctx) error {
	var newNote Note

	if err := c.BodyParser(&newNote); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if newNote.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}

	s.mu.Lock()
	newNote.ID = s.nextID
	s.nextID++
	s.notes = append(s.notes, newNote)
	s.mu.Unlock()

	return c.Status(fiber.StatusCreated).JSON(newNote)
}

func (s *InMemoryStore) UpdateNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var updatedData Note
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, note := range s.notes {
		if note.ID == id {
			s.notes[i].Title = updatedData.Title
			s.notes[i].Content = updatedData.Content
			return c.JSON(s.notes[i])
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note not found"})
}

func (s *InMemoryStore) DeleteNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, note := range s.notes {
		if note.ID == id {
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note not found"})
}
