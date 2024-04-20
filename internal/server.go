package internal
import (
	"fmt"
	"context"
	"database/sql"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	//
	"github.com/bluezy47/Hello-World/internal/services"
	"github.com/bluezy47/Hello-World/internal/oauth"
	"github.com/bluezy47/Hello-World/internal/api"
	//
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/template/html/v2"
)

type Server struct {
	listenAddr string;
	userService *services.UserService;
	sessionStore *session.Store;
}

func NewServer(ctx context.Context, listenAddr string, db *sql.DB) (*Server, error) {
	return &Server{
		listenAddr: listenAddr,
		userService: services.NewUserService(db),
		sessionStore: session.New(),
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
	s.routes(app, s.sessionStore);

	// start the server
	if err := app.Listen(s.listenAddr); err != nil {
		return err;
	}
	return nil;
}
//
func (s *Server) routes(ap *fiber.App, session *session.Store) {
	routes := ap.Group("/");
	//
	routes.Get("/login", func(c fiber.Ctx) error {
		data := map[string]interface{}{}
		return c.Render("login", data);
	}); 
	//
	//
	routes.Get("helloworld", func(c fiber.Ctx) error {
		// check if the user is authenticated...
		sessionInfo, err := session.Get(c);
		if err != nil {
			// todo : check here!
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			});
		}
		isAuthenticated := sessionInfo.Get("isAuthenticated");
		if isAuthenticated == nil || isAuthenticated == false {
			// return c.Redirect().To("/login");
			fmt.Println("::[Server][routes][helloworld] User is not authenticated ::");
		}
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
			"Users": users,
		}
		return c.Render("home", data, "base");
	});
	//
	// OAuth2.0 Redirect Endpoint... `helloworld/user-auth`
	routes.Get("helloworld/user-auth", func(c fiber.Ctx) error {
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
			session, err := session.Get(c);
			if err != nil {
				// todo: check here!...
				return c.Status(http.StatusInternalServerError).SendString(err.Error());
			}
			session.Set("user", userInfo["email"]);
			session.Set("isAuthenticated", true);
			//
			return c.Redirect().To("/helloworld");
		}
		return oauth2.RedirectGoogleLogin(c);
	});
}
