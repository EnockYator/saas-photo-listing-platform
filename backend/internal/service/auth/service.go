package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/config"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       uint
	Email    string
	Password string
}

// Register new user
func Register(db *gorm.DB, cfg *config.Config, input *RegisterInput) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := User{Email: input.Email, Password: string(hash)}
	if err := db.Create(&user).Error; err != nil {
		return "", err
	}

	return generateJWT(cfg, &user)
}

// Login existing user
func Login(db *gorm.DB, cfg *config.Config, input *LoginInput) (string, error) {
	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return generateJWT(cfg, &user)
}

// Generate JWT
func generateJWT(cfg *config.Config, user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(cfg.JWTDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
