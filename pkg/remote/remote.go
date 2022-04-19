package remote

import (
	"errors"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type Remote struct {
	client  *ssh.Client
	config  *ssh.ClientConfig
	address string
}

type RemoteConfig struct {
	Address string
	User    string
	Timeout time.Duration
}

func NewRemote(config RemoteConfig) *Remote {
	return &Remote{
		address: config.Address,
		config: &ssh.ClientConfig{
			User:            config.User,
			Timeout:         config.Timeout,
			HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		},
	}
}

func (r *Remote) ConnectWithPassword(password string) error {
	r.config.Auth = []ssh.AuthMethod{
		ssh.Password(password),
	}

	conn, err := r.dial()
	if err != nil {
		return err
	}

	r.client = conn
	return nil
}

func (r *Remote) ConnectWithKey(pemBytes []byte) error {
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return err
	}
	r.config.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(signer),
	}

	conn, err := r.dial()
	if err != nil {
		return err
	}

	r.client = conn
	return nil
}

func (r *Remote) Connect() error {
	if r.client != nil {
		return nil
	}

	conn, err := r.dial()
	if err != nil {
		return err
	}

	r.client = conn
	return nil
}

func (r *Remote) dial() (*ssh.Client, error) {
	conn, err := ssh.Dial("tcp", r.address, r.config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (r *Remote) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *Remote) NewSession() (*ssh.Session, error) {
	if r.client == nil {
		return nil, errors.New("remote client is not connected")
	}

	session, err := r.client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
