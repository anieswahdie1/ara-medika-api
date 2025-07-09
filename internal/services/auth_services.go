package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/anieswahdie1/ara-medika-api.git/internal/repositories"
	"github.com/anieswahdie1/ara-medika-api.git/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Login(email, password string) (accessToken, refreshToken string, err error)
	Logout(tokenString string, userID uint) error
	RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error)
	InvalidateOtherSessions(userID uint) error
	StoreToken(userID uint, token string) error
}

type authService struct {
	userRepo    repositories.UserRepository
	redisClient *redis.Client
	cfg         *configs.Config
	logger      *logrus.Logger
}

func NewAuthService(userRepo repositories.UserRepository, redisClient *redis.Client, cfg *configs.Config, logger *logrus.Logger) AuthService {
	return &authService{
		userRepo:    userRepo,
		redisClient: redisClient,
		cfg:         cfg,
		logger:      logger,
	}
}

func (s *authService) Login(email, password string) (string, string, error) {
	// Cari user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		s.logger.Errorf("Failed to find user by email: %v", err)
		return "", "", fmt.Errorf("user lookup failed: %w", err)
	}

	// Debug: Cetak hash yang tersimpan (HANYA untuk development)
	s.logger.WithFields(logrus.Fields{
		"email":       email,
		"stored_hash": user.Password,
		"hash_length": len(user.Password),
		"hash_prefix": getHashPrefix(user.Password),
	}).Debug("Password verification debug")

	// Verifikasi password
	if !utils.CheckPassword(password, user.Password) {
		// Tambahkan informasi debug lebih detail
		s.logger.WithFields(logrus.Fields{
			"input_password_length": len(password),
			"stored_hash_prefix":    getHashPrefix(user.Password),
		}).Warn("Password verification failed")

		return "", "", errors.New("invalid credentials")
	}

	// Invalidasi session lain jika ada
	if err := s.InvalidateOtherSessions(user.ID); err != nil {
		s.logger.Warnf("Failed to invalidate other sessions: %v", err)
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(s.cfg, user.ID, user.Email, string(user.Role))
	if err != nil {
		s.logger.Errorf("Failed to generate access token: %v", err)
		return "", "", errors.New("failed to generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(s.cfg, user.ID)
	if err != nil {
		s.logger.Errorf("Failed to generate refresh token: %v", err)
		return "", "", errors.New("failed to generate token")
	}

	// Simpan token di Redis
	if err := s.StoreToken(user.ID, accessToken); err != nil {
		s.logger.Errorf("Failed to store token: %v", err)
		return "", "", errors.New("failed to complete login")
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Logout(tokenString string, userID uint) error {
	// Tambahkan token ke blacklist
	ctx := context.Background()
	expiry := s.cfg.JWTExpire

	// Hitung sisa waktu token
	claims, err := utils.ValidateToken(s.cfg, tokenString)
	if err == nil {
		remaining := time.Until(claims.ExpiresAt.Time)
		if remaining > 0 {
			expiry = remaining
		}
	}

	// Blacklist token
	err = s.redisClient.Set(ctx, tokenString, "blacklisted", expiry).Err()
	if err != nil {
		s.logger.Errorf("Failed to blacklist token: %v", err)
		return errors.New("failed to logout")
	}

	// Hapus dari active tokens
	err = s.redisClient.Del(ctx, fmt.Sprintf("user:%d:token:%s", userID, tokenString)).Err()
	if err != nil {
		s.logger.Warnf("Failed to remove token from active sessions: %v", err)
	}

	return nil
}

func (s *authService) StoreToken(userID uint, token string) error {
	ctx := context.Background()
	expiry := s.cfg.JWTExpire

	// Simpan mapping token -> user
	err := s.redisClient.Set(ctx,
		fmt.Sprintf("user:%d:token:%s", userID, token),
		"active",
		expiry,
	).Err()

	return err
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	// Validasi refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return "", "", errors.New("invalid user ID in token")
	}

	// Dapatkan user dari database
	user, err := s.userRepo.FindByID(uint(userID))
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Generate new tokens
	newAccessToken, err := utils.GenerateToken(s.cfg, user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(s.cfg, user.ID)
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	// Simpan token baru
	if err := s.StoreToken(user.ID, newAccessToken); err != nil {
		return "", "", errors.New("failed to complete refresh")
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) InvalidateOtherSessions(userID uint) error {
	ctx := context.Background()
	pattern := fmt.Sprintf("user:%d:token:*", userID)

	// Dapatkan semua active tokens untuk user ini
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	// Blacklist semua token yang ada
	for _, key := range keys {
		token := strings.TrimPrefix(key, fmt.Sprintf("user:%d:token:", userID))
		s.redisClient.Set(ctx, token, "blacklisted", s.cfg.JWTExpire)
	}

	// Hapus semua active tokens
	_, err = s.redisClient.Del(ctx, keys...).Result()
	return err
}

// Helper function untuk debugging
func getHashPrefix(hash string) string {
	if len(hash) > 8 {
		return hash[:8] + "..."
	}
	return hash
}
