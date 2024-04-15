package internal
import (
	// standard library
	"fmt"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	// custom internal packages
	"github.com/bluezy47/Hello-World/internal/services"
	"github.com/bluezy47/Hello-World/internal/oauth"
	"github.com/bluezy47/Hello-World/internal/api"
	// third party
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"

	// "net/url"
	"net/http"
	// "encoding/json"
	// "encoding/base64"
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
	sessionStore := session.New();
	app := fiber.New(fiber.Config{
		Views: engine,
	});
	// set the session store
	app.Use(sessionStore);
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
	routes := ap.Group("/");

	// Define the routes here...
	routes.Get("/login", func(c *fiber.Ctx) error {
		data := map[string]interface{}{}
		return c.Render("login", data, "base")
	}); // Home Route

	//
	// TODO: Define a function to handle the routes, Do not define the logic here...
	routes.Get("helloworld", func(c *fiber.Ctx) error {
		//
		// check if the user is authenticated...
		session := c.Locals("session").(*session.Session)
		isAuthenticated := session.Get("isAuthenticated");
		if isAuthenticated == nil || isAuthenticated == false {
			return c.Redirect("/login");
		}
		// all ok...
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
			// set the New session for the user.
			session := c.Locals("session").(*session.Session)
			session.Set("user", userInfo["email"]);
			session.Set("isAuthenticated", true);
			//
			// TODO: Remove unnecessary oprations
			// TODO: Add the logics to check the user validity with the database...
			// TODO: Redirect to the Relevent page from here....
			fmt.Println("User Info: ", userInfo);
		}
		return oauth2.RedirectGoogleLogin(c);
	});
}
