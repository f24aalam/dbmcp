package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var credentialsPath string
var secretsPath string
var masterKey []byte

type Credential struct {
	ID       string `json:"id"`
	Database string `json:"database"`
	Name     string `json:"name"`
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	mcpDir := filepath.Join(home, ".mcp")
	err = os.MkdirAll(mcpDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Storage init error: %v\n", err)
		os.Exit(1)
	}

	credentialsPath = filepath.Join(mcpDir, "credentials.json")
	secretsPath = filepath.Join(mcpDir, "secrets.json")

	// Derive a machine-specific master key from the home directory path + hostname.
	// Not perfect security, but protects against casual file reading.
	hostname, _ := os.Hostname()
	seed := home + "|" + hostname
	hash := sha256.Sum256([]byte(seed))
	masterKey = hash[:]

	// Init credentials file
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		data, _ := json.MarshalIndent([]Credential{}, "", "  ")
		os.WriteFile(credentialsPath, data, 0600)
	}

	// Init secrets file
	if _, err := os.Stat(secretsPath); os.IsNotExist(err) {
		data, _ := json.MarshalIndent(map[string]string{}, "", "  ")
		os.WriteFile(secretsPath, data, 0600)
	}
}

// ── Keyring replacement ──────────────────────────────────────────────────────

func secretSet(id, value string) error {
	encrypted, err := encrypt(masterKey, value)
	if err != nil {
		return err
	}
	secrets, err := loadSecrets()
	if err != nil {
		return err
	}
	secrets[id] = encrypted
	return saveSecrets(secrets)
}

func secretGet(id string) (string, error) {
	secrets, err := loadSecrets()
	if err != nil {
		return "", err
	}
	encrypted, ok := secrets[id]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", id)
	}
	return decrypt(masterKey, encrypted)
}

func secretDelete(id string) error {
	secrets, err := loadSecrets()
	if err != nil {
		return err
	}
	delete(secrets, id)
	return saveSecrets(secrets)
}

func loadSecrets() (map[string]string, error) {
	data, err := os.ReadFile(secretsPath)
	if err != nil {
		return nil, err
	}
	var secrets map[string]string
	err = json.Unmarshal(data, &secrets)
	return secrets, err
}

func saveSecrets(secrets map[string]string) error {
	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(secretsPath, data, 0600)
}

// ── AES-GCM encrypt/decrypt ──────────────────────────────────────────────────

func encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, cipherhex string) (string, error) {
	data, err := hex.DecodeString(cipherhex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// ── Public API (same signatures as before) ───────────────────────────────────

func (c Credential) FilterValue() string { return c.Name }
func (c Credential) Title() string       { return c.Name }
func (c Credential) Description() string { return c.ID }

func SaveCredential(
	dbId *string,
	dbName string,
	dbType string,
	dbConnURL string,
) (string, error) {
	var id string
	if dbId != nil {
		id = *dbId
	} else {
		id = fmt.Sprintf("%s-%d", dbType, time.Now().Unix())
	}

	if err := secretSet(id, dbConnURL); err != nil {
		return "", err
	}

	cred := Credential{ID: id, Name: dbName, Database: dbType}
	if err := appendToFile(cred); err != nil {
		return "", err
	}
	return id, nil
}

func GetCredentialById(id string) (string, string, error) {
	creds, err := ListCredentials()
	if err != nil {
		return "", "", fmt.Errorf("error fetching credentials: %w", err)
	}
	for _, cred := range creds {
		if cred.ID == id {
			connURL, err := secretGet(id)
			if err != nil {
				return "", "", fmt.Errorf("error getting connection URL: %w", err)
			}
			return cred.Database, connURL, nil
		}
	}
	return "", "", fmt.Errorf("connection with id %s not found", id)
}

func ListCredentials() ([]Credential, error) {
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}
	var creds []Credential
	err = json.Unmarshal(data, &creds)
	return creds, err
}

func DeleteCredential(id string) error {
	// Remove from secrets
	_ = secretDelete(id) // best-effort

	creds, err := ListCredentials()
	if err != nil {
		return err
	}
	for i, cred := range creds {
		if cred.ID == id {
			creds = append(creds[:i], creds[i+1:]...)
			break
		}
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsPath, data, 0600)
}

func appendToFile(cred Credential) error {
	creds, err := ListCredentials()
	if err != nil {
		return err
	}
	for i, c := range creds {
		if c.ID == cred.ID {
			creds[i] = cred
			data, _ := json.MarshalIndent(creds, "", "  ")
			return os.WriteFile(credentialsPath, data, 0600)
		}
	}
	creds = append(creds, cred)
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsPath, data, 0600)
}
