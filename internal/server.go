package internal
import (
	// standard library
	"fmt"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	// custom internal packages
	"github.com/bluezy47/Hello-World/internal/services"
	// third party
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"net/url"
	"net/http"
	"encoding/json"
	"encoding/base64"
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
	routes := ap.Group("/");

	// Define the routes here...
	routes.Get("", func(c *fiber.Ctx) error {
		data := map[string]interface{}{}
		return c.Render("login", data, "base")
	}); // Home Route

	//
	// TODO: Define a function to handle the routes, Do not define the logic here...
	routes.Get("helloworld", func(c *fiber.Ctx) error {
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

	// OAuth2.0 Redirect Endpoint... `helloworld/user-auth`
	routes.Get("helloworld/user-auth", func(c *fiber.Ctx) error {
		// TODO: gaccess the client_id and client_secret from the environment variables
		// TODO: Don't define all the OAuth2.0 related logic here, define a seperate package for that...
		REDIRECT_URI := "http://localhost:5050/helloworld/user-auth"
		CLIENT_ID := ""
		CLIENT_SECRET := ""

		//
		code := c.Query("code")
		if code != "" {
			// Consent accepted - proceed
			data := url.Values{
				"code":          {code},
				"client_id":     {CLIENT_ID},
				"client_secret": {CLIENT_SECRET},
				"redirect_uri":  {REDIRECT_URI},
				"grant_type":    {"authorization_code"},
			}

			req, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", data)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			defer req.Body.Close()

			var result map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
	
			idToken := result["id_token"].(string)
			
			// get the use info from the id_token
			accessToken := result["access_token"].(string)

			// Now, make a request to the userinfo endpoint to fetch user information
			userInfoReq, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			userInfoReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			userInfoResp, err := http.DefaultClient.Do(userInfoReq)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			defer userInfoResp.Body.Close()

			var userInfo map[string]interface{}
			if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			//
			decodedToken, _ := base64.RawStdEncoding.DecodeString(idToken)
			decodedJSON := make(map[string]interface{})
			json.Unmarshal(decodedToken, &decodedJSON)
			//
			fmt.Println("Decoded Token: ", decodedJSON)

			// TODO: Remove unnecessary oprations
			// TODO: Add the logics to check the user validity with the database...
			// TODO: Redirect to the Relevent page from here....
		}
		redirectURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=%s&redirect_uri=%s&scope=https://www.googleapis.com/auth/userinfo.email&prompt=select_account", CLIENT_ID, REDIRECT_URI)
		return c.Redirect(redirectURL)
	});
}

