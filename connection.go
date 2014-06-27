package main

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"io/ioutil"
	"os"
)

type Connection struct {
	debug  bool
	target *Target
	ssh    *ssh.Client
}

type Command struct {
	Output     string
	ExitStatus int
	Success    bool
	Error      error
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
		fmt.Println(err)
		result = nil
		return
	}

	result = &Connection{
		target: target,
		ssh:    ssh,
	}

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

func (conn *Connection) Exec(command string) *Command {
	debug(command)

	session, err := conn.NewSession()
	exitStatus := 0

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b
	session.Stderr = &b

	err = session.Run(command)

	if err != nil {
		exitErr, ok := err.(*ssh.ExitError)
		if ok {
			exitStatus = exitErr.ExitStatus()
		}
	}

	return &Command{
		Output:     b.String(),
		ExitStatus: exitStatus,
		Success:    err == nil,
		Error:      err,
	}
}

func (conn *Connection) Run(command string) string {
	return conn.Exec(command).Output
}

func (conn *Connection) FileExists(path string) bool {
	return conn.Exec("test -f " + path).Success
}

func (conn *Connection) DirExists(path string) bool {
	return conn.Exec("test -d " + path).Success
}

func (conn *Connection) SymlinkExists(path string) bool {
	return conn.Exec("test -h " + path).Success
}

func (conn *Connection) ProcessExists(pid string) bool {
	return conn.Exec("ps -p " + pid).Success
}

func (conn *Connection) SvnInstalled() bool {
	return conn.Exec("which svn").Success
}

func (conn *Connection) ReadFile(path string) (content string, err error) {
	if !conn.FileExists(path) {
		err = fmt.Errorf("File does not exist: %s", path)
		return
	}

	result := conn.Exec("cat " + path)

	if result.Success {
		content = result.Output
		err = nil
	} else {
		err = fmt.Errorf("Cant read file %s: %s", path, result.Output)
	}

	return content, err
}

func (conn *Connection) CopyFile(src string, dst string) error {
	result := conn.Exec("cp " + src + " " + "dst")

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
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
