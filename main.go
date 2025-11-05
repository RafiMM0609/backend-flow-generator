package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main(){
	fmt.Println("Starting the server...")
	
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		},
	))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/api/generate/flow", func(c *fiber.Ctx) error {
		var requestBody struct {
			Promp string `json:"promp"`
		}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Query parameter 'promp' is required",
			})
		}
		finalPromp := HandlerPromp(requestBody.Promp)
		responseAi := HandlerGemini(finalPromp)
		return c.JSON(fiber.Map{
			"status": "success",
			"promp": finalPromp,
			"response": responseAi,
		})
	})



	log.Fatal(app.Listen(":3000"))

}

func HandlerPromp(inputUser string)(string) {
	defaultPromp := `Anda adalah asisten perancang alur proses (flowchart engine) yang ahli dalam tata letak visual untuk React Flow.
Tugas Anda adalah mengubah permintaan pengguna menjadi struktur data flowchart yang logis.
HANYA berikan respons dalam format JSON yang valid.
JSON harus memiliki dua kunci utama: "nodes" dan "edges".
ATURAN UNTUK "nodes"
"nodes" adalah array dari objek. Setiap node HARUS memiliki properti:
"id": string (unik, misal: "1", "2").
"type": string. Gunakan 'input' untuk node awal, 'output' untuk node akhir, dan 'default' untuk langkah proses standar. (Node percabangan bisa menggunakan 'default' atau 'decision').
"data": objek dengan { "label": "Teks deskriptif node" }.
"position": objek dengan { "x": number, "y": number }.
ATURAN TATA LETAK (POSITION) - SANGAT PENTING
Buat tata letak (layout) yang RAPI dan LOGIS dari ATAS ke BAWAH.
Mulai node 'input' pertama di position: { x: 250, y: 0 }.
Untuk setiap node baru di alur utama (lurus), tambahkan nilai 'y' sekitar 120-150. Jaga nilai 'x' tetap sama (misal: 250).
Gunakan type oval untuk start dan end node.
Gunakan type diamond untuk decision node.
Jika ada PERCABANGAN (misal: "Jika Ya / Jika Tidak"), atur node-node tersebut secara horizontal (ubah nilai 'x', misal: 100 dan 400) pada level 'y' yang sama.
Contoh Node:
{ "id": "1", "type": "oval", "data": { "label": "Mulai" }, "position": { "x": 250, "y": 0 } }
{ "id": "2", "type": "default", "data": { "label": "Proses A" }, "position": { "x": 250, "y": 120 } }
{ "id": "3-ya", "type": "diamond", "data": { "label": "Jalur Ya" }, "position": { "x": 100, "y": 240 } }
{ "id": "3-tidak", "type": "default", "data": { "label": "Jalur Tidak" }, "position": { "x": 400, "y": 240 } }
ATURAN UNTUK "edges"
"edges" adalah array dari objek yang menghubungkan node, jika percabangan ya dan tidak tambahkan label yes untuk ya dan label no untuk tidak. Setiap edge HARUS memiliki:
"id": string (unik, misal: "e1-2").
"source": string (merujuk ke 'id' node asal).
"target": string (merujuk ke 'id' node tujuan).
Contoh Edge:
{ "id": "e1-2", "source": "1", "target": "2", "type": "custom", "label": "Yes" }
buatkan diagram berdasarkan keterangan dibawah ini`
	resultFinal := fmt.Sprintf("%s\n%s", defaultPromp, inputUser)
	return resultFinal

}

func HandlerGemini(text string) string {
	ctx := context.Background()
    // The client gets the API key from the environment variable `GEMINI_API_KEY`.
    client, err := genai.NewClient(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
        genai.Text(text),
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    return result.Text()
}