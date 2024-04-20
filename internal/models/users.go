package models

import (
	"fmt"
	//
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	usersTable = "users"
	usersFields = "id, name, fullname, email, phone, location, birthday, picture, gender, lastonline"
)

type user struct {
	ID int `json:"id"`
	Name string `json:"name"`
	FullName string `json:"fullname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Location string `json:"location"`
	Birthday string `json:"birthday"`
	Picture string `json:"picture"`
	Gender string `json:"gender"`
	LastOnline string `json:"lastonline"`
}

func (u *user) GetScanArgs() []interface{} {
	return []interface{}{
		&u.ID, 
		&u.Name, 
		&u.FullName, 
		&u.Email, 
		&u.Phone, 
		&u.Location, 
		&u.Birthday,
		&u.Picture,
		&u.Gender,
		&u.LastOnline,
	}
}

// 
type UserModel struct {
	mysqlDB *sql.DB;
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{
		mysqlDB: db,
	};
}

//
//
func (m *UserModel) FetchAll() (map[int]interface{}, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", usersFields, usersTable);
	rows, err := m.mysqlDB.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	//
	users := make(map[int]interface{});
	for rows.Next() {
		var u user;
		err := rows.Scan(u.GetScanArgs()...);
		if err != nil {
			return nil, err;
		}
		users[u.ID] = u;
	}
	return users, nil;
}
