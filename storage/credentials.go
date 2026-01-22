package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
)

var credentialsPath string

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

	mpcDir := filepath.Join(home, ".mcp")
	err = os.MkdirAll(mpcDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Storage init error: %v\n", err)
		os.Exit(1)
	}

	credentialsPath = filepath.Join(mpcDir, "credentials.json")

	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		emptyCredentials := []Credential{}
		data, err := json.MarshalIndent(emptyCredentials, "", " ")

		if err != nil {
			fmt.Fprintf(os.Stderr, "Storage init error: %v\n", err)
			os.Exit(1)
		}

		os.WriteFile(credentialsPath, data, 0644)
	}
}

func (c Credential) FilterValue() string {
	return c.Name
}

func (c Credential) Title() string {
	return c.Name
}

func (c Credential) Description() string {
	return c.ID
}

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

	err := keyring.Set("mcp-db-connections", id, dbConnURL)
	if err != nil {
		return "", err
	}

	cred := Credential{
		ID:       id,
		Name:     dbName,
		Database: dbType,
	}

	err = appendToFile(cred)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetCredentialById(id string) (string, string, error) {
	creds, err := ListCredentials()
	if err != nil {
		return "", "", fmt.Errorf("Error fetching credentials: %w", err)
	}

	for _, cred := range creds {
		if cred.ID == id {
			connUrl, err := keyring.Get("mcp-db-connections", id)

			if err != nil {
				return "", "", fmt.Errorf("Error while getting connection URL: %w", err)
			}

			return cred.Database, connUrl, nil
		}
	}

	return "", "", fmt.Errorf("Connection with id %s not found", id)
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

func appendToFile(cred Credential) error {
	creds, err := ListCredentials()

	if err != nil {
		return err
	}

	for i, c := range creds {
		if c.ID == cred.ID {
			creds[i] = cred
			data, _ := json.MarshalIndent(creds, "", " ")

			return os.WriteFile(credentialsPath, data, 0644)
		}
	}

	// Apend new
	creds = append(creds, cred)
	data, err := json.MarshalIndent(creds, "", " ")

	if err != nil {
		return err
	}

	return os.WriteFile(credentialsPath, data, 0644)
}

func DeleteCredential(id string) error {
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

	data, err := json.MarshalIndent(creds, "", " ")

	if err != nil {
		return err
	}

	return os.WriteFile(credentialsPath, data, 0644)
}
