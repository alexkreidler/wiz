package local

import (
	"github.com/alexkreidler/wiz/executor"
	"os"
	"os/exec"
	"strconv"
)

type Environment struct {
	Port int64
}

func (e Environment) Name() string {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return "local environment: " + host
}

// No configuration necessary for the local environment

// Maybe able to configure the port that the executor listens on
func (e Environment) Configure(d interface{}) error {
	//panic("implement me")
	//e = d.(*Environment)
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
	c := exec.Command("wiz", "executor", "--port", strconv.FormatInt(e.Port, 10))

	c.Stdout = os.Stdout
	return c.Run()
}

func NewEnvironment() executor.Environment {
	return Environment{Port: 9003}
}

func init() {
	err := executor.RegisterEnvironment("local", NewEnvironment)
	if err != nil {
		panic(err)
	} // .(executor.EnvironmentConstructor))
}
