package utils

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengubah password plaintext menjadi hash menggunakan bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword memverifikasi password plaintext dengan hash yang tersimpan
func CheckPassword(password, hashedPassword string) bool {
	// Validasi input
	if len(hashedPassword) == 0 {
		logrus.Warn("Empty hashed password provided")
		return false
	}

	// Bandingkan password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("Password comparison failed")
		return false
	}

	return true
}
