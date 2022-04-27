package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type SFTP struct {
	name string
	basedir string
	host string
	username string
	password string
	privatekey string
	timeout time.Duration
	keyexchanges []string
}

func NewSFTP(name string, cfg map[string]string) SFTP {
	c := SFTP{
		name: name,
		basedir: cfg["basedir"],
		host: cfg["host"],
		username: cfg["username"],
		password: cfg["password"],
		keyexchanges: []string{"diffie-hellman-group-exchange-sha256", "diffie-hellman-group14-sha256"},
		timeout: time.Second * 30,
	}
	return c
}


func (c *SFTP) getClient() (*sftp.Client, error) {
	var auths []ssh.AuthMethod
	if c.password != "" {
        auths = append(auths, ssh.Password(c.password))
    }
	if c.privatekey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(c.privatekey))
		if err != nil {
			return nil, err
		}
		auths = append(auths, ssh.PublicKeys(signer))
	}
 	cfg := &ssh.ClientConfig{
		User: c.username,
		Auth: auths,
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error { return nil },
		Timeout:         c.timeout,
		Config: ssh.Config{
			KeyExchanges: c.keyexchanges,
		},
	}
	sshClient, err := ssh.Dial("tcp", c.host, cfg)
	if err != nil {
		return nil, err
	}
 	return sftp.NewClient(sshClient)
}

func (s SFTP) Name() string {
	return s.name
}

// Get returns a file (if exist) from the storage target
func (s SFTP) Get(filename string) ([]byte, error) {
	sftpClient, err := s.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect [%v]", err)
	}
	defer sftpClient.Close()
	dest := strings.Join([]string{s.basedir, filename}, "/")
	destFile, err := sftpClient.OpenFile(dest, (os.O_RDONLY))
    if err != nil {
        return nil, fmt.Errorf("failed to open {%s} [%v]", dest, err)
    }
    defer destFile.Close()
	dd := new(bytes.Buffer)
	wr := bufio.NewWriter(dd)
	n, err := destFile.WriteTo(wr)
	if err != nil {
		return nil, err
	}
	err = wr.Flush()
	if err != nil {
		return nil, err
	}
	log.Debugf("transferred %v bytes", n)
	if n == 0 {
		return nil, fmt.Errorf("failed to transfer bytes")
	}
	return dd.Bytes(), nil
}

// Put saves a file to the storage target
func (s SFTP) Put(filename string, data []byte) error {
	sftpClient, err := s.getClient()
	if err != nil {
		return fmt.Errorf("failed to connect [%v]", err)
	}
	defer sftpClient.Close()
	dest := strings.Join([]string{s.basedir, filename}, "/")
	destFile, err := sftpClient.OpenFile(dest, (os.O_WRONLY|os.O_CREATE|os.O_TRUNC))
    if err != nil {
        return fmt.Errorf("failed to create {%s} [%v]", dest, err)
    }
    defer destFile.Close()
	n, err := destFile.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("uploaded %v of %v", n, len(data))
	}
	return nil
}

func (s SFTP) Download(name string, filename string) (string, error) {
	sftpClient, err := s.getClient()
	if err != nil {
		return "", fmt.Errorf("failed to connect [%v]", err)
	}
	defer sftpClient.Close()
	sftpFile := strings.Join([]string{s.basedir, name}, "/")
	destFile, err := sftpClient.OpenFile(sftpFile, (os.O_RDONLY))
    if err != nil {
        return "", fmt.Errorf("failed to open {%s} [%v]", sftpFile, err)
    }
    defer destFile.Close()
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sh := enigma.NewStreamHash(destFile)
	n, err := io.Copy(f, sh.GetReader())
	if err != nil || n == 0 {
		return "", fmt.Errorf("transferred %v bytes -> %v", n, err)
	}
	return sh.GetHash(), nil
}

func (s SFTP) Upload(filename string, name string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	sftpClient, err := s.getClient()
	if err != nil {
		return fmt.Errorf("failed to connect [%v]", err)
	}
	defer sftpClient.Close()
	sftpFile := strings.Join([]string{s.basedir, name}, "/")
	destFile, err := sftpClient.OpenFile(sftpFile, (os.O_WRONLY|os.O_CREATE|os.O_TRUNC))
    if err != nil {
        return fmt.Errorf("failed to create {%s} [%v]", sftpFile, err)
    }
    defer destFile.Close()
	n, err := destFile.ReadFrom(r)
	if err != nil {
		return err
	}
	if n != info.Size() {
		return fmt.Errorf("uploaded %v of %v", n, info.Size())
	}
	return nil
}