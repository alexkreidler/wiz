package local

import (
	"fmt"
	"github.com/alexkreidler/wiz/environment"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
)

type Environment struct {
	Port uint32
}

func (e Environment) Name() string {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return "local environment: " + host
}

// No configuration necessary for the local environment

func checkPortAvailable(port uint32) bool {
	_, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(port), 10))

	if err != nil {
		return false
	}

	return true
}

// Maybe able to configure the port that the executor listens on
func (e *Environment) Configure(d interface{}) error {
	//panic("implement me")
	newEnv, ok := d.(Environment)
	if !ok {
		return fmt.Errorf("failed to convert configuration to environment")
	}

	if !checkPortAvailable(newEnv.Port) {
		return fmt.Errorf("invalid port: in use")
	}

	*e = newEnv

	return nil
}

func (e Environment) GetCurrentConfiguration() interface{} {
	//panic("implement me")
	return nil
}

func (e Environment) IsValidConfiguration() bool {
	// check if port is in valid port range
	return true
}

func (e Environment) StartExecutor(node string) error {
	c := exec.Command("wiz", "executor", "--port", ":"+strconv.FormatUint(uint64(e.Port), 10))
	err := c.Start()
	if err != nil {
		log.Println("failed with err", err)
		return err
	}
	log.Println("started command", c)

	return nil
}

func (e Environment) Describe() environment.SerializableEnv {
	host, _ := os.Hostname()
	return environment.SerializableEnv{
		EnvironmentID: "local",
		Description:   "Local Environment (" + host + ")",
		Host:          "localhost:" + strconv.FormatUint(uint64(e.Port), 10),
		Configuration: e,
	}
}

func NewEnvironment() environment.Environment {
	return &Environment{Port: 9003}
}

//
//func init() {
//	err := environment.RegisterEnvironment("local", NewEnvironment)
//	if err != nil {
//		panic(err)
//	} // .(executor.EnvironmentConstructor))
//}
