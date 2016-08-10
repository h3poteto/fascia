package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type deploy struct {
	User string
	Host string
	Port string

	DockerImageName     string
	DockerImageTag      string
	FirstContainerName  string
	SecondContainerName string
	FirstContainerPort  int
	SecondContainerPort int

	SharedDirectory string

	HostName string

	client *ssh.Client
}

// 新しいdocker imageをpullする
// 現在動いているdockerのポートを確認する
// 新しいdockerを起動する
// db:migrate
// confdがwatchしているredisのキーを変更し，新しいコンテナのホストとポートを送る
// confdのwatch interval+バッファ分だけだけ待つ
// curlしてみて通信できることを確認する
// 古いdockerをstopする
// 80番にcurlしてみて通信できることを確認する
// 古いdockerコンテナを削除する
// 古いdocker imageを削除する

func main() {
	d, err := initialize()
	if err != nil {
		panic(err)
	}
	err = d.prepareDockerImage()
	if err != nil {
		panic(err)
	}

	port, err := d.checkRunningPort()
	if err != nil {
		panic(err)
	}

	var newPort int
	var newContainer, oldContainer string
	switch port {
	case d.FirstContainerPort:
		newPort = d.SecondContainerPort
		newContainer = d.SecondContainerName
		oldContainer = d.FirstContainerName
	case d.SecondContainerPort:
		newPort = d.FirstContainerPort
		newContainer = d.FirstContainerName
		oldContainer = d.SecondContainerName
	case 0:
		// どちらのコンテナも起動していないパターン
		newPort = d.FirstContainerPort
		newContainer = d.FirstContainerName
	default:
		panic("Container is running unexpected port")
	}

	err = d.removeOldContainer()
	// 起動中のコンテナを消そうとした場合にはエラーが帰ってくるが，特に問題はないので表示するだけ
	if err != nil {
		log.Println(err)
	}

	// migration
	err = d.migration()
	if err != nil {
		panic(err)
	}

	err = d.startNewContainer(newContainer, newPort)
	if err != nil {
		panic(err)
	}

	err = d.refreshRedis(newPort)
	if err != nil {
		panic(err)
	}
	// confdのrefresh intervalを10[s]にしているので，10[s]以上待つ必要がある
	time.Sleep(20 * time.Second)

	// 一応curlが通ることを確認してから進めたい
	_, err = d.checkServiceLiving()
	if err != nil {
		panic(err)
	}

	if len(oldContainer) > 0 {
		err = d.stopOldContainer(oldContainer)
		if err != nil {
			panic(err)
		}
		// 古いコンテナを停止した段階で，もう一度curlしたい
		_, err = d.checkServiceLiving()
		if err != nil {
			panic(err)
		}
	}

	err = d.removeOldContainer()
	if err != nil {
		log.Println(err)
	}

	err = d.removeOldImages()
	if err != nil {
		log.Println(err)
	}
	log.Println("Deploy success!")
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
		User:                m["user"].(string),
		Host:                m["host"].(string),
		Port:                m["port"].(string),
		DockerImageName:     m["docker_image_name"].(string),
		DockerImageTag:      m["docker_image_tag"].(string),
		FirstContainerName:  m["first_container_name"].(string),
		SecondContainerName: m["second_container_name"].(string),
		FirstContainerPort:  m["first_container_port"].(int),
		SecondContainerPort: m["second_container_port"].(int),
		SharedDirectory:     m["shared_directory"].(string),
		HostName:            m["host_name"].(string),
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

func (d *deploy) prepareDockerImage() error {
	session := d.getSession()
	defer session.Close()

	log.Println("Docker pull")
	command := fmt.Sprintf("docker pull %v:%v", d.DockerImageName, d.DockerImageTag)
	log.Println(command)
	if err := session.Run(command); err != nil {
		return err
	}
	return nil
}

func (d *deploy) checkRunningPort() (int, error) {

	firstContainerPort, _ := d.runningPort(d.FirstContainerName, d.FirstContainerPort)
	secondContainerPort, _ := d.runningPort(d.SecondContainerName, d.SecondContainerPort)

	// 両方コンテナが起動している場合は想定外ななのでエラーにする
	if firstContainerPort != 0 && secondContainerPort != 0 {
		return 0, errors.New("Both containers are running")
	} else if firstContainerPort != 0 {
		return firstContainerPort, nil
	} else if secondContainerPort != 0 {
		return secondContainerPort, nil
	}

	// 両方共起動していない場合は，あらたに起動すればいいだけなので，エラーにはしない
	return 0, nil
}

func (d *deploy) runningPort(name string, reservedPort int) (int, error) {
	session := d.getSession()
	defer session.Close()

	session.Stdout = nil
	session.Stderr = nil

	log.Println("Check running docker port")
	command := fmt.Sprintf("docker port %v", name)
	log.Println(command)
	result, err := session.CombinedOutput(command)
	log.Println(string(result))
	if err != nil {
		return 0, err
	}
	if len(result) <= 0 {
		return 0, errors.New("Can not find docker port")
	}
	port, err := pickupPortNumber(string(result), reservedPort)
	if err != nil {
		return 0, err
	}
	return port, nil
}

func pickupPortNumber(dockerPort string, reservedPort int) (int, error) {
	regEx := `0.0.0.0:` + strconv.Itoa(reservedPort)
	r := regexp.MustCompile(regEx)
	match := r.FindAllStringSubmatch(dockerPort, -1)
	if len(match) != 1 {
		return 0, errors.New("Cannot find port")
	}
	if len(match[0]) != 1 {
		return 0, errors.New("Cannot find port")
	}
	return reservedPort, nil
}

func (d *deploy) startNewContainer(name string, port int) error {
	session := d.getSession()
	defer session.Close()

	log.Println("Start new container")
	command := fmt.Sprintf("docker run -d -v %v:/root/fascia/public/statics --env-file /home/ubuntu/.docker-env --name %v -p %v:9090 %v", d.SharedDirectory, name, port, d.DockerImageName)
	log.Println(command)
	if err := session.Run(command); err != nil {
		return err
	}
	return nil
}

func (d *deploy) migration() error {
	session := d.getSession()
	defer session.Close()

	log.Println("db migration")
	command := fmt.Sprintf("docker run --rm --env-file /home/ubuntu/.docker-env %v gom exec goose -env production up", d.DockerImageName)
	if err := session.Run(command); err != nil {
		return err
	}
	return nil
}

func (d *deploy) refreshRedis(port int) error {
	session := d.getSession()
	defer session.Close()

	command := fmt.Sprintf("bash -l -c 'redis-cli -h $REDIS_HOST -p $REDIS_PORT set /app/upstream 127.0.0.1:%v'", port)
	log.Println(command)
	if err := session.Run(command); err != nil {
		return err
	}
	return nil
}

func (d *deploy) stopOldContainer(name string) error {
	session := d.getSession()
	defer session.Close()

	log.Println("Stop old container")
	command := fmt.Sprintf("docker stop %v", name)
	log.Println(command)
	if err := session.Run(command); err != nil {
		return err
	}
	return nil
}

func (d *deploy) checkServiceLiving() (int, error) {
	session := d.getSession()
	defer session.Close()

	session.Stdout = nil
	session.Stderr = nil

	log.Println("Check service is living")
	command := fmt.Sprintf("curl --insecure -H '%v' https://127.0.0.1 -o /dev/null -w '%%{http_code}' -s", d.HostName)
	log.Println(command)
	result, err := session.CombinedOutput(command)
	if err != nil {
		return 0, err
	}
	if len(result) <= 0 {
		return 0, errors.New("Can not get HTTP status code")
	}
	statusCode, err := strconv.Atoi(string(result))
	if err != nil {
		return 0, err
	}
	if statusCode != 200 {
		return statusCode, errors.New("HTTP status code is not 200")
	}
	return statusCode, nil
}

func (d *deploy) removeOldContainer() error {
	session := d.getSession()
	defer session.Close()

	log.Println("Remove old docker container")
	command := "docker rm `docker ps -a -q`"
	log.Println(command)
	if err := session.Run(command); err != nil {
		return err
	}

	return nil
}

func (d *deploy) removeOldImages() error {
	session := d.getSession()
	defer session.Close()

	log.Println("Remove old docker images")
	command := fmt.Sprintf("docker rmi -f $(docker images | awk '/<none>/ { print $3 }')")
	log.Println(command)
	if err := session.Run(command); err != nil {
		return err
	}
	return nil
}
