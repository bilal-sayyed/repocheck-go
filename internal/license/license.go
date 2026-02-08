package license

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type License struct {
	Plan    string `json:"plan"`
	Variant string `json:"variant"`
	Exp     int64  `json:"exp"`
}

// Replace with real public key
// Public Key (Ed25519) - ROTATED Jan 2026
var publicKey = ed25519.PublicKey{
	0x6e, 0xb6, 0xec, 0xdc, 0x85, 0x1, 0x7e, 0x0,
	0xe8, 0x52, 0x6d, 0x2a, 0x5b, 0xd7, 0xfd, 0x11,
	0xaf, 0x6c, 0x29, 0x3f, 0x2f, 0x20, 0xba, 0xa7,
	0x9f, 0x95, 0xc2, 0x9, 0x9a, 0x1c, 0x4d, 0x7,
}

func GetLicensePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".repocheck_license"
	}
	return filepath.Join(home, ".repocheck_license")
}

func SaveLicense(key string) error {
	return os.WriteFile(GetLicensePath(), []byte(strings.TrimSpace(key)), 0600)
}

func LoadLicense() (string, error) {
	data, err := os.ReadFile(GetLicensePath())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func VerifyLicense(key string) (*License, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 3 || parts[0] != "RC1" {
		return nil, errors.New("invalid license format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid signature encoding")
	}

	if !ed25519.Verify(publicKey, payload, signature) {
		return nil, errors.New("invalid license signature")
	}

	var lic License
	if err := json.Unmarshal(payload, &lic); err != nil {
		return nil, errors.New("invalid license payload")
	}

	if time.Now().Unix() > lic.Exp {
		return nil, errors.New("license expired")
	}

	return &lic, nil
}

func IsPro() bool {
	key, err := LoadLicense()
	if err != nil || key == "" {
		return false
	}

	lic, err := VerifyLicense(key)
	if err != nil {
		return false
	}

	return lic.Plan == "pro"
}
