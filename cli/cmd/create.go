/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/drawiin/go-expert/cli/internal/commands"
	"github.com/drawiin/go-expert/cli/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	// Init dependencies
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	categoryDB := database.NewCategoryDB(db)

	// Configuring command
	createCmd := newCreateCmd(categoryDB, func() {
		db.Close()
	})
	categoryCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("name", "n", "", "Name of the category")
	createCmd.Flags().StringP("description", "d", "Empry description", "Description of the category")
	createCmd.Flags().BoolP("showInfo", "i", false, "Should log the info")
	createCmd.MarkFlagsRequiredTogether("name", "description")
}

func newCreateCmd(db *database.CategoryDB, closeDb func()) *cobra.Command {
	return &cobra.Command{
		Use:      "create",
		Short:    "Create a new category",
		RunE:     createRunCommand(db),
		PostRunE: createPostRunCommand(closeDb),
	}
}

func createRunCommand(db *database.CategoryDB) commands.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		showInfo, _ := cmd.Flags().GetBool("showInfo")
		category, err := db.Create(name, description)
		if err != nil {
			return err
		}
		println("Category created")
		if showInfo {
			fmt.Println("Id ", category.ID)
			fmt.Println("Name: ", name)
			fmt.Println("Description: ", description)
		}
		return nil
	}
}

func createPostRunCommand(shutDown func()) commands.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		defer shutDown()
		return nil
	}
}
