package main

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Connection struct {
	debug  bool
	target *Target
	ssh    *ssh.Client
}

func NewConnection(target *Target) (result *Connection, err error) {
	privateKey := parsekey(privateKeyPath())

	config := &ssh.ClientConfig{
		User: target.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
			ssh.Password(target.password),
		},
	}

	ssh, err := ssh.Dial("tcp", target.getAddress(), config)

	if err != nil {
		result = nil
		return
	}

	result = &Connection{
		target: target,
		ssh:    ssh,
	}

	result.debug = true

	return
}

func (conn *Connection) NewSession() (session *ssh.Session, err error) {
	session, err = conn.ssh.NewSession()

	if err != nil {
		session = nil
		return
	}

	termModes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err = session.RequestPty("xterm", 80, 40, termModes); err != nil {
		session = nil
		return
	}

	return
}

func (conn *Connection) Run(command string) string {
	session, err := conn.NewSession()
	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b
	session.Stderr = &b

	if conn.debug {
		log.Println(command)
	}

	err = session.Run(command)

	if err != nil {
		fmt.Println("Failed to run: ", err.Error())
	}

	return b.String()
}

func privateKeyPath() string {
	return fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
}

func parsekey(file string) ssh.Signer {
	privateBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	return private
}
