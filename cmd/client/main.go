// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"golang.org/x/crypto/ssh"
// 	"golang.org/x/term"
// )

// var (
// 	user     = flag.String("l", "", "login_name")
// 	password = flag.String("pass", "", "password")
// 	port     = flag.Int("p", 22, "port")
// )

// func main() {
// 	flag.Parse()
// 	if flag.NArg() == 0 {
// 		flag.Usage()
// 		os.Exit(2)
// 	}

// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
// 	ctx, cancel := context.WithCancel(context.Background())

// 	go func() {
// 		if err := run(ctx); err != nil {
// 			log.Print(err)
// 		}
// 		cancel()
// 	}()

// 	select {
// 	case <-sig:
// 		cancel()
// 	case <-ctx.Done():
// 	}
// }

// func run(ctx context.Context) error {
// 	fmt.Println("Connecting...", *user, *password, *port)
// 	config := &ssh.ClientConfig{
// 		User: *user,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(*password),
// 		},
// 		Timeout:         5 * time.Second,
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	hostport := fmt.Sprintf("%s:%d", flag.Arg(0), *port)
// 	conn, err := ssh.Dial("tcp", hostport, config)
// 	if err != nil {
// 		return fmt.Errorf("cannot connect %v: %v", hostport, err)
// 	}
// 	defer conn.Close()

// 	session, err := conn.NewSession()
// 	if err != nil {
// 		return fmt.Errorf("cannot open new session: %v", err)
// 	}
// 	defer session.Close()

// 	go func() {
// 		<-ctx.Done()
// 		conn.Close()
// 	}()

// 	fd := int(os.Stdin.Fd())
// 	state, err := term.MakeRaw(fd)
// 	if err != nil {
// 		return fmt.Errorf("terminal make raw: %s", err)
// 	}
// 	defer term.Restore(fd, state)

// 	w, h, err := term.GetSize(fd)
// 	if err != nil {
// 		return fmt.Errorf("terminal get size: %s", err)
// 	}

// 	modes := ssh.TerminalModes{
// 		ssh.ECHO:          1,
// 		ssh.TTY_OP_ISPEED: 14400,
// 		ssh.TTY_OP_OSPEED: 14400,
// 	}

// 	term := os.Getenv("TERM")
// 	if term == "" {
// 		term = "xterm-256color"
// 	}
// 	if err := session.RequestPty(term, h, w, modes); err != nil {
// 		return fmt.Errorf("session xterm: %s", err)
// 	}

// 	session.Stdout = os.Stdout
// 	session.Stderr = os.Stderr
// 	session.Stdin = os.Stdin

// 	if err := session.Shell(); err != nil {
// 		return fmt.Errorf("session shell: %s", err)
// 	}

// 	if err := session.Wait(); err != nil {
// 		if e, ok := err.(*ssh.ExitError); ok {
// 			switch e.ExitStatus() {
// 			case 130:
// 				return nil
// 			}
// 		}
// 		return fmt.Errorf("ssh: %s", err)
// 	}
// 	return nil
// }

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
