package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type deploy struct {
	User string
	Host string
	Port string

	Repository string
	Branch     string
	Name       string
	SourceFile string

	DeployTo string

	client *ssh.Client
}

func main() {
	d, err := initialize()
	if err != nil {
		panic(err)
	}
	d.checkParentDirectory()
	d.createSrcDirectory()
	d.createBinDirectory()
	d.cloneSourceCode()
	d.build()
	d.migration()
	d.restart()
}

func initialize() (d *deploy, e error) {
	root := os.Getenv("GOJIROOT")
	path := filepath.Join(root, "deploy/setting.yml")
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}
	d = &deploy{
		User:       m["user"].(string),
		Host:       m["host"].(string),
		Port:       m["port"].(string),
		Repository: m["repository"].(string),
		Branch:     m["branch"].(string),
		Name:       m["name"].(string),
		SourceFile: m["sourcefile"].(string),
		DeployTo:   m["deployto"].(string),
	}
	return d, nil
}

func (d *deploy) initClient() {
	user := d.User

	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		panic(err)
	}
	signers, err := agent.NewClient(sock).Signers()
	if err != nil {
		panic(err)
	}
	auths := []ssh.AuthMethod{ssh.PublicKeys(signers...)}
	config := &ssh.ClientConfig{
		User: user,
		Auth: auths,
	}
	d.client, err = ssh.Dial("tcp", d.Host+d.Port, config)
	if err != nil {
		panic(err)
	}
}

func (d *deploy) getSession() *ssh.Session {
	if d.client == nil {
		d.initClient()
	}

	session, err := d.client.NewSession()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err != nil {
		panic(err)
	}

	return session
}

func (d *deploy) checkParentDirectory() {
	session := d.getSession()
	defer session.Close()

	log.Printf("Check directory deploy to: %v", d.DeployTo)
	command := fmt.Sprintf("if ! [ -d %v ]; then echo 'cloud not found %v'; fi", d.DeployTo, d.DeployTo)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

func (d *deploy) createSrcDirectory() {
	session := d.getSession()
	defer session.Close()

	log.Printf("Create src directoy in deploy to: %v", d.DeployTo)
	command := fmt.Sprintf("if ! [ -d %v/src ]; then mkdir -p %v/src; fi", d.DeployTo, d.DeployTo)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

func (d *deploy) createBinDirectory() {
	session := d.getSession()
	defer session.Close()

	log.Printf("Create bin directoy in deploy to: %v", d.DeployTo)
	command := fmt.Sprintf("if ! [ -d %v/bin ]; then mkdir -p %v/bin; fi", d.DeployTo, d.DeployTo)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

func (d *deploy) cloneSourceCode() {
	session := d.getSession()
	defer session.Close()

	log.Printf("Clone source codes from github: %v", d.Repository)
	command := fmt.Sprintf("if ! [ -d %v/src/%v ]; then cd %v/src && git clone -b %v %v; else cd %v/src/%v && git pull origin %v; fi", d.DeployTo, d.Name, d.DeployTo, d.Branch, d.Repository, d.DeployTo, d.Name, d.Branch)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

func (d *deploy) setEnvironments() {
	session := d.getSession()
	defer session.Close()
}

func (d *deploy) beforeBuild() {
	session := d.getSession()
	defer session.Close()

	log.Print("Prepare build")
	command := fmt.Sprintf("source $HOME/.bash_profile; cd %v/src/%v && gom install", d.DeployTo, d.Name)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

func (d *deploy) build() {
	d.beforeBuild()
	session := d.getSession()
	defer func() {
		session.Close()
		d.afterBuild()
	}()

	log.Print("Build go source")
	command := fmt.Sprintf("source $HOME/.bash_profile; cd %v/src/%v && gom build -o %v/bin/%v %v", d.DeployTo, d.Name, d.DeployTo, d.Name, d.SourceFile)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

func (d *deploy) afterBuild() {
	session := d.getSession()
	defer session.Close()

	log.Println("After build")
	command := fmt.Sprintf("source $HOME/.bash_profile; cd %v/src/%v; npm install; npm run-script js-release; npm run-script sass-release", d.DeployTo, d.Name)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

// db migrate
func (d *deploy) migration() {
	session := d.getSession()
	defer session.Close()

	log.Println("db migration")
	command := fmt.Sprintf("source $HOME/.bash_profile; cd %v/src/%v && gom exec goose -env production up", d.DeployTo, d.Name)
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}

func (d *deploy) restart() {
	session := d.getSession()
	defer session.Close()

	log.Println("restart application")
	command := fmt.Sprint("sudo service supervisor restart")
	log.Println(command)
	if err := session.Run(command); err != nil {
		panic(err)
	}
}
