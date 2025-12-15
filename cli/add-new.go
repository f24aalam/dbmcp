package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/f24aalam/godbmcp/database"
	"github.com/f24aalam/godbmcp/storage"
)

func AddNewConnection(cred *storage.Credential) {
	var dbId *string
	var dbName string
	var dbType string
	var dbConnUrl string

	if cred != nil {
		dbId = &cred.ID
		dbName = cred.Name
		dbType = cred.Database
		_, dbConnUrl, _ = storage.GetCredentialById(cred.ID)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Connection Name").
				Validate(huh.ValidateNotEmpty()).
				Value(&dbName),
			huh.NewSelect[string]().
				Title("Select Database").
				Options(
					huh.NewOption("MySQL", "mysql"),
				).
				Validate(huh.ValidateNotEmpty()).
				Value(&dbType),
			huh.NewInput().
				Title("Enter connection string").
				Validate(huh.ValidateNotEmpty()).
				Value(&dbConnUrl),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	conn := database.Connection{
		Database: dbType,
		ConnectionUrl: dbConnUrl,
	}

	err = conn.Open()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	id, err := storage.SaveCredential(dbId, dbName, dbType, dbConnUrl)
	if err != nil {
		fmt.Println("Error in saving connection: ", err)
		return
	}

	fmt.Println("Database connection success, saved with id: ", id)

	defer conn.Close()
}
