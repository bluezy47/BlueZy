package main
import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)


func main() {
	// make the HTML engine
	engine := html.New("./templates", ".html");
	app := fiber.New(fiber.Config{
		Views: engine,
	});

	// serve the static files
	app.Static("/static", "./templates/static")

	// define the routes
	app.Get("/", func(c *fiber.Ctx) error {
		data := map[string]interface{}{
			"Title": "Hello, World!",
			"Message": "This Message is from thea Server",
		}
		//
		// RENDER THE HOME PAGE
		return c.Render("home", data, "base");
	})
	//
	// start the server
	log.Println("Server is running on http://localhost:3000");
	log.Fatal(app.Listen(":3000"));
}
