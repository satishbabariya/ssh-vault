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
	"context"
	"fmt"
	"os"
	"ssh-vault/internal/proto"
	"time"

	"github.com/cli/oauth"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/square/go-jose.v2/jwt"
)

func init() {
	godotenv.Load()
}

type VaultClient struct {
	conn   *grpc.ClientConn
	client proto.AuthServiceClient
}

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
				ctx, cancel := context.WithTimeout(c.Context, 30*time.Second)
				defer cancel()

				conn, err := grpc.Dial(
					"localhost:1203",
					grpc.WithInsecure(),
					grpc.WithBlock(),
					// grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor),
					// grpc.WithStreamInterceptor(interceptor.StreamClientInterceptor),
				)
				if err != nil {
					return err
				}

				client := proto.NewAuthServiceClient(conn)

				vault := &VaultClient{
					conn:   conn,
					client: client,
				}

				token, err := keyring.Get("vault", "token")
				if err != nil {
					t, err := vault.Login(ctx)
					if err != nil {
						return err
					}
					token = *t
				}

				t, err := jwt.ParseSigned(token)
				if err != nil {
					return err
				}

				var claims jwt.Claims

				err = t.UnsafeClaimsWithoutVerification(&claims)
				if err != nil {
					return err
				}

				// check if token is expired
				if claims.Expiry.Time().Before(time.Now()) {
					t, err := vault.Login(ctx)
					if err != nil {
						return err
					}
					token = *t
				}

				fmt.Println(token)

				return conn.Close()
			},
		},
		{
			Name:  "list",
			Usage: `List all SSH Vault remote hosts`,
			Action: func(c *cli.Context) error {

				interceptor := ClientInterceptor{}

				conn, err := grpc.Dial(
					"localhost:1203",
					grpc.WithInsecure(),
					grpc.WithBlock(),
					grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor),
				)
				if err != nil {
					return err
				}

				return conn.Close()
			},
		},
	}

	// Run the app.
	if err := app.Run(os.Args); err != nil {
		// Log the error and exit.
		logrus.Errorln(err)
	}

}

func (v *VaultClient) Login(ctx context.Context) (*string, error) {
	t, err := v.AuthenticateWithGithub(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to store token: %v", err)
	}

	err = keyring.Set(
		"vault", "token", *t,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to store token: %v", err)
	}

	return t, nil
}

func (v *VaultClient) AuthenticateWithGithub(ctx context.Context) (*string, error) {
	config, err := v.client.GetConfig(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}

	if config.GithubClientId == "" {
		return nil, fmt.Errorf("github client id is empty")
	}

	flow := &oauth.Flow{
		Host:     oauth.GitHubHost(config.GithubHost),
		ClientID: config.GithubClientId,
		Scopes: []string{
			"user:email",
		},
	}

	accessToken, err := flow.DeviceFlow()
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Authenticate(ctx, &proto.AuthenticateRequest{
		Token: accessToken.Token,
	})

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("authenticate response is nil")
	}

	return &resp.Token, nil
}

// ClientInterceptor is a gRPC interceptor that adds the access token to the request
type ClientInterceptor struct {
	accessToken string
}

// UnaryClientInterceptor is a gRPC interceptor that adds the access token to the request
func (interceptor *ClientInterceptor) UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	md.Set("authorization", interceptor.accessToken)

	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
