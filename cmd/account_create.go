package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/server"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// accountCreateCmd represents the accountCreate command
var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new blockchain account",
	Long:  `create a new blockchain account`,
	Run: func(cmd *cobra.Command, args []string) {

		var password1, password2 string
		passwordMatch := false

		fmt.Println("Password is required for account creation. It can't be restored so please save it in a safe place.")
		fmt.Println("Please, enter password: \n")

		for !passwordMatch {
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal("Fatal error while creating an account.")
			}
			password1 = string(bytePassword)

			if strings.Trim(password1, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please, enter passwords again.")
				continue
			}

			fmt.Println("Enter password again: ")
			bytePassword, err = term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(err)
			}

			password2 = string(bytePassword)

			if password1 == password2 {
				passwordMatch = true
			} else {
				fmt.Println("Passwords do not match. Please, enter passwords again.")
			}

		}
		accountStr, err := account.CreateAccount(password1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Your new account's address is:", accountStr)

		server.Start(accountStr, "48658")

	},
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
}
