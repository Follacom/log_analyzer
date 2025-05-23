/*
Copyright © 2025 Mathieu FOLLACO mathieu.follaco@hotmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"log_analyzer/lib"

	"github.com/gofiber/fiber/v2"
)

func main() {
	lib.InitConfig()

	lib.ListenAccess()
	// lib.ListenError()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// Route: Get processed logs
	// app.Get("/logs", func(c *fiber.Ctx) error {
	// 	lib.MuDatabase.Lock()
	// 	defer lib.MuDatabase.Unlock()
	// 	db, err := lib.RotateDatabase()
	// 	if err != nil {
	// 		lib.LogError(err)
	// 		return c.Send([]byte(err.Error()))
	// 	}
	// 	var data []model.ApacheAccessLog
	// 	_ = db.Find(&data)
	// 	res := make([]interface{}, 0)
	// 	for _, row := range data {
	// 		res = append(res, row)
	// 	}
	// 	return c.JSON(len(res))
	// })

	// Start server on port 8080
	app.Listen(":8080")
}
