package internal
import (
	// standard library
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	// custom internal packages
	"github.com/bluezy47/Hello-World/internal/services"
	// third party
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Server struct {
	listenAddr string;
	userService *services.UserService;
}

func NewServer(ctx context.Context, listenAddr string, db *sql.DB) (*Server, error) {
	return &Server{
		listenAddr: listenAddr,
		userService: services.NewUserService(db),
	}, nil
}
//
func (s *Server) Start() error {
	engine := html.New("../templates", ".html");
	app := fiber.New(fiber.Config{
		Views: engine,
	});
	// serve the static files
	app.Static("/static", "../templates/static")
	s.routes(app);
	//
	// start the server
	if err := app.Listen(s.listenAddr); err != nil {
		return err;
	}
	return nil;
}
//
func (s *Server) routes(ap *fiber.App) {
	routes := ap.Group("/helloworld");
	//
	// TODO: Define a function to handle the routes, Do not define the logic here...
	routes.Get("/", func(c *fiber.Ctx) error {
		//
		users, err := s.userService.FetchAll();
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			});
		}
		//
		data := map[string]interface{}{
			"Title": "Hello, World!",
			"Message": "This Message is from the Server",
			"Users": users, // pass the users to the view
		}
		return c.Render("home", data, "base");
	});
}