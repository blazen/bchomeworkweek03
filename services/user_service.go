package services

import (
	"errors"

	"gorm.io/gorm"

	"gin-examples/project/models"
	"gin-examples/project/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, utils.NewAppError(409, "Username already exists")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, utils.NewAppError(409, "Email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashedPassword),
		PostNumber: 0,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(404, "User not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(401, "Invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, utils.NewAppError(401, "Invalid credentials")
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 如果更新邮箱，检查是否已存在
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, utils.NewAppError(409, "Email already exists")
		}
		user.Email = req.Email
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// 查询每个用户审计合法、非法的文章数
func (s *UserService) StatisticPostAuditStatus() ([]models.StatisticUserResponse, error) {
	var sta []models.StatisticUserResponse
	if err := s.db.Raw(`
		SELECT 
			u.username,
			p.audit_status,
			count(*) AS AuditStatusCount
		FROM posts p
		LEFT JOIN users u on p.user_id = u.id 
		GROUP BY u.username, p.audit_status
	`).Scan(&sta).Error; err != nil {
		return nil, err
	}
	return sta, nil
}
