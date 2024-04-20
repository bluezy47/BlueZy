package services
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	// 
	"github.com/bluezy47/Hello-World/internal/models"
)

type UserService struct {
	userModel *models.UserModel;
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		userModel: models.NewUserModel(db),
	};
}

func (s *UserService) FetchAll() (map[int]interface{}, error) {
	users, err := s.userModel.FetchAll();
	if err != nil {
		return nil, err;
	}
	//
	return users, nil;
}
