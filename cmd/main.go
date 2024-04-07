package main
import (
	"os"
	"fmt"
	"context"
	//
	"github.com/bluezy47/Hello-World/internal"
	"github.com/bluezy47/Hello-World/pkg/sql"
)


func main() {
	ctx := context.Background();

	// TODO: Intergrate the Redis Connection...

	// init the DB connection
	dbConfig := sql.DBConfig{
		Username: "root",
		Password: "12345678",
		Host: "localhost",
		Port: "3306",
		Database: "bluezy_chat",
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
	}
	dbConn, err := sql.ConnInit(dbConfig);
	if err != nil {
		fmt.Println("[SQL] Conncetion Initilization Failed!", err);
		os.Exit(1);
	}
	fmt.Println("[SQL] Connection Initilized Successfully!");

	//
	helloworldServer, err := internal.NewServer(ctx, "127.0.0.1:5050", dbConn);
	if err != nil {
		fmt.Println("Error: ", err);
		return;
	}
	//
	helloworldServer.Start();
	fmt.Println("Server is running ... :5050");
}
