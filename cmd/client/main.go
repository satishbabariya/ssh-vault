package main

import (
	"os"

	"github.com/satishbabariya/vault/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

func main() {

	app := &cli.App{
		Name:        "sshv",
		Usage:       "SSH Vault",
		Description: `SSH Vault is a tool for managing SSH keys.`,
	}

	app.Commands = []*cli.Command{
		{
			Name:  "login",
			Usage: `Login to SSH Vault`,
			Action: func(c *cli.Context) error {
				vault, err := client.NewClientUnsafe(c.Context, "http://localhost:1203")
				if err != nil {
					return err
				}

				err = vault.Login(c.Context)
				if err != nil {
					return err
				}

				logrus.Info("Login successful")

				return nil
			},
		},
		{
			Name:  "logout",
			Usage: `Logout from SSH Vault`,
			Action: func(c *cli.Context) error {
				err := keyring.Delete("vault", "token")
				if err != nil {
					return err
				}

				logrus.Info("Logout successful")

				return nil
			},
		},
		{
			Name:  "list",
			Usage: `List all SSH Vault remote hosts`,
			Action: func(c *cli.Context) error {
				vault, err := client.NewClient(c.Context, "http://localhost:1203")
				if err != nil {
					return err
				}

				logrus.Info("Listing SSH Vault remote hosts", vault)

				return nil
			},
		},
	}

	// Run the app.
	if err := app.Run(os.Args); err != nil {
		// Log the error and exit.
		logrus.Errorln(err)
	}
}
