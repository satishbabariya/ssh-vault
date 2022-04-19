package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/satishbabariya/vault/pkg/remote"
	"github.com/satishbabariya/vault/pkg/store"
	"github.com/satishbabariya/vault/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func main() {
	// db filepath
	dbFile := "vault.db"
	dbfilepath := fmt.Sprintf("%s/.vault/%s", os.Getenv("HOME"), dbFile)

	// create directory if not exists
	if _, err := os.Stat(fmt.Sprintf("%s/.vault", os.Getenv("HOME"))); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s/.vault", os.Getenv("HOME")), 0755)
	}

	store, err := store.Open(dbfilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	app := &cli.App{
		Name:        "vault",
		Usage:       "SSH Vault",
		Description: `SSH Vault is a tool for managing SSH keys.`,
	}

	app.Commands = []*cli.Command{
		{
			Name:  "list",
			Usage: `List all SSH Vault remote hosts`,
			Action: func(c *cli.Context) error {

				remotes, err := store.Credentials()
				if err != nil {
					return err
				}

				if len(remotes) == 0 {
					logrus.Info("No remote hosts found")
					return nil
				}

				// TODO: print remotes
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.SetStyle(table.StyleLight)
				t.AppendHeader(table.Row{
					"Host",
					"Port",
					"Tags",
				})
				t.SetColumnConfigs([]table.ColumnConfig{
					{Number: 1, Align: text.AlignRight},
					{Number: 2, Align: text.AlignCenter},
					{Number: 3, Align: text.AlignCenter},
				})

				for _, remote := range remotes {
					t.AppendRow(table.Row{
						remote.Host,
						remote.Port,
						remote.TagsString(),
					})
				}

				t.Render()

				return nil
			},
		},
		{
			Name:  "add",
			Usage: `Add a new SSH Vault remote host`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "host",
					// Aliases:  []string{"h"},
					Usage:    "Hostname of the remote host",
					Required: true,
				},
				&cli.IntFlag{
					Name: "port",
					// Aliases: []string{"p"},
					Usage: "Port of the remote host",
					Value: 22,
				},
				&cli.StringFlag{
					Name: "user",
					// Aliases:  []string{"u"},
					Usage:    "Username of the remote host",
					Required: true,
				},
				&cli.StringFlag{
					Name: "private-key",
					// Aliases:   []string{"k"},
					Usage:     "Path to the private key of the remote host",
					TakesFile: true,
				},
				&cli.StringFlag{
					Name: "password",
					// Aliases: []string{"pw"},
					Usage: "Password of the remote host",
				},
				&cli.StringSliceFlag{
					Name:    "tags",
					Aliases: []string{"t"},
					Usage:   "Tags of the remote host",
				},
			},
			Action: func(c *cli.Context) error {
				host := c.String("host")
				port := c.Int("port")
				user := c.String("user")
				privateKey := c.String("private-key")
				password := c.String("password")
				tags := c.StringSlice("tags")

				credential := &types.Credential{
					Host: host,
					Port: port,
					User: user,
					Tags: tags,
				}

				if privateKey != "" {
					pemBytes, err := ioutil.ReadFile(privateKey)
					if err != nil {
						return err
					}
					credential.PrivatKey = pemBytes
				}

				if password != "" {
					credential.Password = &password
				}

				err := store.Add(credential)
				if err != nil {
					return err
				}

				logrus.Infof("Added remote host %s", host)

				return nil
			},
		},
		{
			Name:  "connect",
			Usage: `Connect to a SSH Vault remote host`,
			Action: func(c *cli.Context) error {

				// get the first argument
				host := c.Args().First()
				if host == "" {
					return errors.New("host or tag name is required")
				}

				// get the credentials
				credentials, err := store.Credentials()
				if err != nil {
					return err
				}

				// find the credentials
				var credential *types.Credential
				for _, c := range credentials {
					if c.Host == host {
						credential = &c
						break
					}

					if c.HasTag(host) {
						credential = &c
						break
					}
				}

				if credential == nil {
					return errors.New("credential not found")
				}

				sig := make(chan os.Signal, 1)
				signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
				ctx, cancel := context.WithCancel(c.Context)

				go func() {
					if err := run(ctx, *credential); err != nil {
						logrus.Error(err)
					}
					cancel()
				}()

				select {
				case <-sig:
					cancel()
				case <-ctx.Done():
				}

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

func run(ctx context.Context, credential types.Credential) error {
	// connect to the remote host
	conn := remote.NewRemote(remote.RemoteConfig{
		Address: fmt.Sprintf("%s:%d", credential.Host, credential.Port),
		User:    credential.User,
	})

	if len(credential.PrivatKey) > 0 {
		err := conn.ConnectWithKey(credential.PrivatKey)
		if err != nil {
			return err
		}
	} else if credential.Password != nil {
		err := conn.ConnectWithPassword(*credential.Password)
		if err != nil {
			return err
		}
	} else {
		err := conn.Connect()
		if err != nil {
			return err
		}
	}

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %v", err)
	}
	defer session.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("terminal make raw: %s", err)
	}
	defer term.Restore(fd, state)

	w, h, err := term.GetSize(fd)
	if err != nil {
		return fmt.Errorf("terminal get size: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}
	if err := session.RequestPty(term, h, w, modes); err != nil {
		return fmt.Errorf("session xterm: %s", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("session shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return fmt.Errorf("ssh: %s", err)
	}
	return nil
}
