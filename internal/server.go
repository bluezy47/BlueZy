package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	//
	"github.com/bluezy47/Hello-World/internal/api"
	"github.com/bluezy47/Hello-World/internal/oauth"
	"github.com/bluezy47/Hello-World/internal/services"

	//
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
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
	//
	// set the middleware for the ws 
	// todo: please move this to a separate function or package in future...
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	});
	//
	// serve the static files
	app.Static("/static", "../templates/static")
	s.routes(app);


	// start the server
	if err := app.Listen(s.listenAddr); err != nil {
		return err;
	}
	return nil;
}
//
func (s *Server) routes(ap *fiber.App) {
	routes := ap.Group("/");
	//
	// user login page
	routes.Get("/login", func(c *fiber.Ctx) error {
		data := map[string]interface{}{}
		return c.Render("login", data);
	}); 
	//
	// Home Page
	routes.Get("helloworld", func(c *fiber.Ctx) error {
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
			"Users": users,
		}
		return c.Render("home", data, "base");
	});
	//
	// Handle the websocket connection...
	routes.Get("ws/helloworld", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))
		log.Println(c.Params("id"))
		log.Println(c.Query("email"))
		log.Println(c.Cookies("session"))
		// todo: here you can add some authentication logic if you want...
		// call the websocket handler
		services.WebsocketHandler(c);
	}));
	//
	// OAuth2.0 Redirect Endpoint... `helloworld/user-auth`
	routes.Get("helloworld/user-auth", func(c *fiber.Ctx) error {
		oauth2, err := oauth.NewOAuth();
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		//
		code := c.Query("code")
		if code != "" {
			result, err := oauth2.HandleInitCode(code);
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error());
			}
			accessToken := result["access_token"].(string)
			userInfo, err := api.GetGoogleUserDetails(accessToken);
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error());
			}
			fmt.Println(userInfo); // todo: in future, use this data to make a session user!.
			return c.Redirect("/helloworld");
		}
		return oauth2.RedirectGoogleLogin(c);
	});
}
