package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	app := fiber.New(fiber.Config{
		ServerHeader: "Jacques' Autobahnapi",
	})

	app.Get("/:value", func(c *fiber.Ctx) error {
		fmt.Println(c.Params("value"))

		res, err := client.Get(context.Background(), (c.Params("value"))).Result()
		if err != nil {
			return c.SendString(fmt.Sprint(err))
		} else {
			fmt.Println("Success getting value in Redis")

		}

		fmt.Println(res)
		return c.SendString(res)

	})

	app.Listen(":3000")
}
