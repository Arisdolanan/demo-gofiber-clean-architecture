package usecase

import (
	"errors"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/messaging/kafka"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/redis"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AuthUsecase interface {
	Register(email, password string, schoolID *int64, userType entity.UserType) error
	Login(email, password string) (*entity.AuthToken, error)
	RefreshToken(refreshToken string) (*entity.AuthToken, error)
	Logout(userID int64, accessToken string) error
	VerifyToken(token string) (*entity.User, error)
	// Enhanced security methods
	ValidatePasswordComplexity(password string) error
	BlacklistToken(token string, expiration time.Duration) error
}

type authUsecase struct {
	authRepo      postgresql.AuthRepository
	authRedis     redis.AuthRedisRepository
	emailUsecase  EmailUsecase
	validate      *validator.Validate
	log           *logrus.Logger
	jwtSecret     string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	kafkaProducer *kafka.UserProducer
}

func NewAuthUsecase(
	authRepo postgresql.AuthRepository,
	authRedis redis.AuthRedisRepository,
	emailUsecase EmailUsecase,
	validate *validator.Validate,
	log *logrus.Logger,
	jwtSecret string,
	kafkaProducer *kafka.UserProducer,
) AuthUsecase {
	return &authUsecase{
		authRepo:      authRepo,
		authRedis:     authRedis,
		emailUsecase:  emailUsecase,
		validate:      validate,
		log:           log,
		jwtSecret:     jwtSecret,
		accessTTL:     15 * time.Minute,
		refreshTTL:    7 * 24 * time.Hour,
		kafkaProducer: kafkaProducer,
	}
}

// Register handles user registration
func (uc *authUsecase) Register(email, password string, schoolID *int64, userType entity.UserType) error {
	// Validate password complexity first
	if err := uc.ValidatePasswordComplexity(password); err != nil {
		uc.log.Errorf("Password complexity validation failed for %s: %v", email, err)
		return err
	}

	// Check if user already exists
	existingUser, err := uc.authRepo.FindByEmail(email)
	if err != nil {
		uc.log.Errorf("Error checking existing user: %v", err)
		return err
	}

	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Hash password with complexity validation
	hashedPassword, err := utils.ValidateAndHashPassword(password)
	if err != nil {
		uc.log.Errorf("Error validating and hashing password: %v", err)
		return err
	}

	// Create user
	user := &entity.User{
		Email:       email,
		Password:    hashedPassword,
		SchoolID:    schoolID,
		UserType:    userType,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.authRepo.Register(user); err != nil {
		uc.log.Errorf("Error registering user: %v", err)
		return err
	}

	// Send verification email
	if err := uc.emailUsecase.SendVerificationEmail(user.ID, email); err != nil {
		uc.log.Errorf("Error sending verification email: %v", err)
		// Continue with registration even if email fails
	}

	uc.log.Infof("User registered successfully: %s", email)
	return nil
}

// Login handles user authentication
func (uc *authUsecase) Login(email, password string) (*entity.AuthToken, error) {
	// Find user by email
	user, err := uc.authRepo.FindByEmail(email)
	if err != nil {
		uc.log.Errorf("Error finding user: %v", err)
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, uc.jwtSecret, uc.accessTTL)
	if err != nil {
		uc.log.Errorf("Error generating access token: %v", err)
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Email, uc.jwtSecret, uc.refreshTTL)
	if err != nil {
		uc.log.Errorf("Error generating refresh token: %v", err)
		return nil, err
	}

	// Store refresh token in Redis
	if err := uc.authRedis.StoreRefreshToken(user.ID, refreshToken, uc.refreshTTL); err != nil {
		uc.log.Errorf("Error storing refresh token: %v", err)
		return nil, err
	}

	authToken := &entity.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(uc.accessTTL),
	}

	// Send login event to Kafka
	if uc.kafkaProducer != nil {
		uc.log.Infof("Publishing user login event for user: %s", email)
		if err := uc.kafkaProducer.PublishUserEvent(user); err != nil {
			uc.log.Errorf("Error publishing user login event: %v", err)
		}
	} else {
		uc.log.Infof("Kafka producer is disabled, skipping user login event")
	}

	uc.log.Infof("User logged in successfully: %s", email)
	return authToken, nil
}

// VerifyToken validates a token and returns the associated user
func (uc *authUsecase) VerifyToken(tokenString string) (*entity.User, error) {
	// Validate token
	claims, err := utils.ValidateToken(tokenString, uc.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Find user by ID
	user, err := uc.authRepo.FindByEmail(claims.Email)
	if err != nil {
		uc.log.Errorf("Error finding user: %v", err)
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// RefreshToken generates new tokens using a refresh token
func (uc *authUsecase) RefreshToken(refreshToken string) (*entity.AuthToken, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken, uc.jwtSecret)
	if err != nil {
		uc.log.Errorf("Invalid refresh token: %v", err)
		return nil, errors.New("invalid refresh token")
	}

	// Verify refresh token exists in Redis and matches
	storedToken, err := uc.authRedis.GetRefreshToken(claims.UserID)
	if err != nil {
		uc.log.Errorf("Error getting refresh token from Redis: %v", err)
		return nil, errors.New("invalid refresh token")
	}

	if storedToken != refreshToken {
		uc.log.Warnf("Refresh token mismatch for user ID: %d", claims.UserID)
		return nil, errors.New("invalid refresh token")
	}

	// Delete used refresh token (one-time use for security)
	if err := uc.authRedis.DeleteRefreshToken(claims.UserID); err != nil {
		uc.log.Errorf("Error deleting refresh token: %v", err)
		return nil, err
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(claims.UserID, claims.Email, uc.jwtSecret, uc.accessTTL)
	if err != nil {
		uc.log.Errorf("Error generating access token: %v", err)
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := utils.GenerateToken(claims.UserID, claims.Email, uc.jwtSecret, uc.refreshTTL)
	if err != nil {
		uc.log.Errorf("Error generating refresh token: %v", err)
		return nil, err
	}

	// Store new refresh token in Redis
	if err := uc.authRedis.StoreRefreshToken(claims.UserID, newRefreshToken, uc.refreshTTL); err != nil {
		uc.log.Errorf("Error storing refresh token: %v", err)
		return nil, err
	}

	authToken := &entity.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(uc.accessTTL),
	}

	uc.log.Infof("Token refreshed successfully for user: %s", claims.Email)
	return authToken, nil
}

// Logout handles user logout
func (uc *authUsecase) Logout(userID int64, accessToken string) error {
	// Blacklist the current access token
	if err := uc.authRedis.BlacklistToken(accessToken, uc.accessTTL); err != nil {
		uc.log.Errorf("Error blacklisting access token during logout: %v", err)
		return err
	}

	// Delete refresh token from Redis
	if err := uc.authRedis.DeleteRefreshToken(userID); err != nil {
		uc.log.Errorf("Error deleting refresh token during logout: %v", err)
		return err
	}

	uc.log.Infof("User logged out successfully: %d", userID)
	return nil
}

// ValidatePasswordComplexity validates password against complexity requirements
func (uc *authUsecase) ValidatePasswordComplexity(password string) error {
	return utils.ValidatePasswordComplexity(password)
}

// BlacklistToken adds a token to the blacklist
func (uc *authUsecase) BlacklistToken(token string, expiration time.Duration) error {
	return uc.authRedis.BlacklistToken(token, expiration)
}
