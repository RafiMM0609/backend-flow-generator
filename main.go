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

	app.Get("/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString("Hello, " + name + "!")
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



	log.Fatal(app.Listen(":5000"))

}

func HandlerPromp(inputUser string)(string) {
	defaultPromp := `Anda adalah asisten perancang visual yang ahli dalam membuat struktur data Flowchart dan Diagram Hubungan Entitas (ERD) untuk React Flow.
Tugas Anda adalah menganalisis permintaan pengguna.
Jika permintaan adalah tentang alur proses/langkah-langkah, hasilkan Flowchart JSON.
Jika permintaan adalah tentang struktur database/tabel, hasilkan ERD JSON.
HANYA berikan respons dalam format JSON yang valid.
JSON harus memiliki dua kunci utama: "nodes" dan "edges".

1. Aturan untuk Flowchart
Penerapan: Jika permintaan berkaitan dengan alur proses/langkah-langkah.
ATURAN UNTUK "nodes" (Flowchart)
"nodes" adalah array dari objek. Setiap node HARUS memiliki properti:
"id": string (unik, misal: "1", "2").
"type": string. Gunakan oval untuk node mulai dan akhir. Gunakan diamond untuk node keputusan/percabangan. Gunakan default untuk langkah proses standar.
"data": objek dengan { "label": "Teks deskriptif node" }.
"position": objek dengan { "x": number, "y": number }.

ATURAN TATA LETAK (POSITION) - SANGAT PENTING (Flowchart)
Buat tata letak (layout) yang RAPI dan LOGIS dari ATAS ke BAWAH.
Mulai node pertama di position: { x: 250, y: 0 }.
Untuk setiap node baru di alur utama (lurus), tambahkan nilai 'y' sekitar 120-150. Jaga nilai 'x' tetap sama (misal: 250).
Jika ada PERCABANGAN (misal: "Jika Ya / Jika Tidak"), atur node-node hasil percabangan secara horizontal (ubah nilai 'x', misal: 100 dan 400) pada level 'y' yang sama.
ATURAN UNTUK "edges" (Flowchart)
"edges" adalah array dari objek yang menghubungkan node. Setiap edge HARUS memiliki:
"id": string (unik, misal: "e1-2").
"source": string (merujuk ke 'id' node asal).
"target": string (merujuk ke 'id' node tujuan).
Tambahkan properti "label" (misal: "Yes", "No") untuk edge yang keluar dari node keputusan (diamond).
2. Aturan untuk ERD
Penerapan: Jika permintaan berkaitan dengan struktur database/tabel.
ATURAN UNTUK "nodes" (ERD) - Menggunakan Template Anda
"nodes" adalah array dari objek. Setiap node HARUS mewakili satu tabel/entitas dan memiliki properti:
"id": string (unik, misal: "erd-1", "erd-users").
"type": HARUS menggunakan tableNode.
"position": objek dengan { "x": number, "y": number }. Terapkan tata letak yang rapi, tempatkan tabel secara terpisah dan mudah dibaca, umumnya dengan jarak minimal x: 300 dan y: 300 antar tabel.
"data": objek yang berisi:
"tableName": string (nama tabel, misal: "users").
"columns": array kolom/atribut. Setiap kolom harus memiliki:
"id": string (unik, misal: "col-u1").
"name": string (nama kolom).
"type": string (tipe data, misal: "INT", "VARCHAR").
"isPK": boolean (true jika Primary Key, false jika bukan).
"isFK": boolean (true jika Foreign Key, false jika bukan).
ATURAN UNTUK "edges" (ERD) - Menggunakan Template Anda
"edges" adalah array dari objek yang merepresentasikan hubungan antar tabel. Setiap edge HARUS memiliki:
"id": string (unik, misal: "erd-e1").
"source": string (merujuk ke 'id' node tabel asal).
"sourceHandle": string. Formatnya: [id_kolom_PK]-src (misal: "col-u1-src").
"target": string (merujuk ke 'id' node tabel tujuan).
"targetHandle": string. Formatnya: [id_kolom_FK]-tgt (misal: "col-p2-tgt").
"type": HARUS menggunakan relationship.
"data": objek dengan { "relationship": "Kardinalitas hubungan" } (misal: "1:N", "1:1", "N:M").
Contoh Respons (Jika Pengguna Meminta Flowchart):
JSON
{
  "nodes": [
    { "id": "1", "type": "oval", "data": { "label": "Mulai" }, "position": { "x": 250, "y": 0 } },
    { "id": "2", "type": "default", "data": { "label": "Lakukan Login" }, "position": { "x": 250, "y": 120 } },
    { "id": "3", "type": "diamond", "data": { "label": "Login Berhasil?" }, "position": { "x": 250, "y": 240 } },
    { "id": "4-ya", "type": "default", "data": { "label": "Tampilkan Dashboard" }, "position": { "x": 100, "y": 360 } },
    { "id": "4-tidak", "type": "default", "data": { "label": "Tampilkan Pesan Error" }, "position": { "x": 400, "y": 360 } },
    { "id": "5", "type": "oval", "data": { "label": "Selesai" }, "position": { "x": 250, "y": 480 } }
  ],
  "edges": [
    { "id": "e1-2", "source": "1", "target": "2" },
    { "id": "e2-3", "source": "2", "target": "3" },
    { "id": "e3-4ya", "source": "3", "target": "4-ya", "label": "Yes" },
    { "id": "e3-4tidak", "source": "3", "target": "4-tidak", "label": "No" },
    { "id": "e4ya-5", "source": "4-ya", "target": "5" },
    { "id": "e4tidak-5", "source": "4-tidak", "target": "5" }
  ]
}
Contoh Respons (Jika Pengguna Meminta ERD untuk Users dan Posts):
JSON
{
  "nodes": [
    {
      "id": "erd-users",
      "type": "tableNode",
      "position": { "x": 100, "y": 100 },
      "data": {
        "tableName": "users",
        "columns": [
          { "id": "col-u1", "name": "id", "type": "INT", "isPK": true, "isFK": false },
          { "id": "col-u2", "name": "username", "type": "VARCHAR(50)", "isPK": false, "isFK": false },
          { "id": "col-u3", "name": "email", "type": "VARCHAR(100)", "isPK": false, "isFK": false }
        ]
      }
    },
    {
      "id": "erd-posts",
      "type": "tableNode",
      "position": { "x": 500, "y": 100 },
      "data": {
        "tableName": "posts",
        "columns": [
          { "id": "col-p1", "name": "id", "type": "INT", "isPK": true, "isFK": false },
          { "id": "col-p2", "name": "user_id", "type": "INT", "isPK": false, "isFK": true },
          { "id": "col-p3", "name": "title", "type": "VARCHAR(255)", "isPK": false, "isFK": false },
          { "id": "col-p4", "name": "content", "type": "TEXT", "isPK": false, "isFK": false }
        ]
      }
    }
  ],
  "edges": [
    {
      "id": "e-users-posts",
      "source": "erd-users",
      "sourceHandle": "col-u1-src",
      "target": "erd-posts",
      "targetHandle": "col-p2-tgt",
      "type": "relationship",
      "data": { "relationship": "1:N" }
    }
  ]
} 
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