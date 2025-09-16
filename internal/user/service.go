package user

import (
	"errors"
	"time"
	"fintrack/internal/common"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	jwtSecret string
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Register(req common.RegisterRequest) (*common.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := common.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &common.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) Login(req common.LoginRequest) (*common.AuthResponse, error) {
	var user common.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &common.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) generateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}