package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"strconv"
)

type Todo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{
	{ID: 1, Name: "julio", Completed: false},
	{ID: 2, Name: "maria", Completed: true},
}

func getTodo(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(todos)
}

func postTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name      string `json:"name"`
		Completed bool   `json:"completed"`
	}
	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	particularTodo := Todo{
		ID:        len(todos) + 1,
		Name:      body.Name,
		Completed: body.Completed,
	}
	todos = append(todos, particularTodo)
	return ctx.Status(fiber.StatusCreated).JSON(particularTodo)

}

func deleteTodo(ctx *fiber.Ctx) error {
	paramID := ctx.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Id",
		})
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[0:i], todos[i+i:]...)
			return ctx.Status(fiber.StatusOK).JSON(todo)
		}
	}
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Record not found",
	})
}

func patchTodo(ctx *fiber.Ctx) error {

	type request struct {
		Name      string `json:"name"`
		Completed bool   `json:"completed"`
	}
	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parser JSON",
		})
	}

	paramID := ctx.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Id",
		})
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos[i] = Todo{
				ID:        id,
				Name:      body.Name,
				Completed: body.Completed,
			}
			return ctx.Status(fiber.StatusOK).JSON(todos[i])
		}
	}
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Record not found",
	})
}

func main() {

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Welcome\n")
	})

	app.Get("/todo", getTodo)
	app.Post("/todo", postTodo)
	app.Delete("/todo/:id", deleteTodo)
	app.Patch("/todo/:id", patchTodo)

	app.Listen(":3000")
}
