package sql
import (
	"fmt"
	"log"
	"sync"
	//
	// mysql driver...
	"database/sql"
	"github.com/go-sql-driver/mysql"
)
//
// package level variables
var (
	sqlDB *sql.DB
	dbOnce sync.Once // `sync.Once` is used to make sure that the connection is created only once and is thread safe.
)
//
type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	MaxOpenConnections int
	MaxIdleConnections int
}
//
func ConnInit(cfg DBConfig) (*sql.DB, error) {
	configMap := map[string]string {
		"username": cfg.Username,
		"password": cfg.Password,
		"host": cfg.Host,
		"port": cfg.Port,
		"database": cfg.Database,
	}
	// set values to `MaxOpenConnections` and `MaxIdleConnections`
	if cfg.MaxOpenConnections == 0 {
		cfg.MaxOpenConnections = 10
	}
	if cfg.MaxIdleConnections == 0 {
		cfg.MaxIdleConnections = 5
	}
	//
	// validate the config
	for k, v := range configMap {
		if v == "" {
			return nil, fmt.Errorf("missing %s in config", k)
		}
	}
	//
	dbOnce.Do(func() {
		config := mysql.Config{
			User:                 cfg.Username,
			Passwd:               cfg.Password,
			Net:                  "tcp",
			Addr:                 cfg.Host + ":" + cfg.Port,
			DBName:               cfg.Database,
			AllowNativePasswords: true, // allows the use of the native password authentication method
			CheckConnLiveness:    true, // check liveness of connection before using it
			MultiStatements: 	  true, // allow multiple statements in one query
		}
		var err error;
		sqlDB, err = sql.Open("mysql", config.FormatDSN())
		if err != nil {
			log.Fatal("[MySQL] DB Connection Error: ", err);
			return;
		}
		//
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
		log.Println("[MySQL] DB Connection Established");
	});
	//
	return sqlDB, nil;
}
