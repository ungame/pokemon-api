package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"net/http"
	"pokemon-api/assets"
	"strconv"
)

var port int

func init() {
	flag.IntVar(&port, "p", 8080, "set api port")
	assets.Load()
	flag.Parse()
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})

	app.Get("/pokemons", func(ctx *fiber.Ctx) error {
		query := ctx.Query("limit", "10")
		limit, err := strconv.Atoi(query)
		if err != nil || limit < 10 || limit > assets.Count() {
			limit = 10
		}
		names := assets.GetNames()[:limit]
		return ctx.Status(http.StatusOK).JSON(names)
	})

	app.Get("/pokemons/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")
		pokemon := assets.Get(name)
		if len(pokemon) == 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}
		return ctx.Status(http.StatusOK).JSON(Pokemon{
			Name: name,
			Gif:  pokemon,
		})
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}

type Pokemon struct {
	Name string `json:"name"`
	Gif  []byte `json:"gif"`
}
