package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
}

func main() {
	// Читаємо змінні оточення, які передасть Docker Compose.
	// Якщо запускаємо локально без Докера, підставляться дефолтні значення (fallback).
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "admin")
	dbPassword := getEnv("DB_PASSWORD", "1234")
	dbName := getEnv("DB_NAME", "contacts_db")

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Неможливо підключитися до БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}

	appContext := &App{DB: db}
	app := fiber.New()

	app.Get("/contacts", appContext.GetAllContacts)
	app.Get("/contacts/:id", appContext.GetContactByID)
	app.Post("/contacts", appContext.CreateContact)
	app.Put("/contacts/:id", appContext.UpdateContact)
	app.Delete("/contacts/:id", appContext.DeleteContact)

	fmt.Println("Сервер стартує на порту 8080...")
	log.Fatal(app.Listen(":8080"))
}

// Помічник для читання Env-змінних
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (a *App) AlbaniaContacts(c *fiber.Ctx) {}

func (a *App) GetAllContacts(c *fiber.Ctx) error {
	rows, err := a.DB.Query("SELECT id, name, phone FROM contacts ORDER BY id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.Name, &contact.Phone); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		contacts = append(contacts, contact)
	}

	if contacts == nil {
		return c.JSON([]Contact{})
	}

	return c.JSON(contacts)
}

func (a *App) GetContactByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var contact Contact
	err = a.DB.QueryRow("SELECT id, name, phone FROM contacts WHERE id = $1", id).
		Scan(&contact.ID, &contact.Name, &contact.Phone)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Contact not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(contact)
}

func (a *App) CreateContact(c *fiber.Ctx) error {
	var contact Contact
	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if contact.Name == "" || contact.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and phone are required"})
	}

	query := "INSERT INTO contacts (name, phone) VALUES ($1, $2) RETURNING id"
	err := a.DB.QueryRow(query, contact.Name, contact.Phone).Scan(&contact.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(contact)
}

func (a *App) UpdateContact(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var contact Contact
	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	query := "UPDATE contacts SET name = $1, phone = $2 WHERE id = $3"
	res, err := a.DB.Exec(query, contact.Name, contact.Phone, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Contact not found"})
	}

	contact.ID = id
	return c.JSON(contact)
}

func (a *App) DeleteContact(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	query := "DELETE FROM contacts WHERE id = $1"
	res, err := a.DB.Exec(query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Contact not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
