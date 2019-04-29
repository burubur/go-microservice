package cmd

import (
	"log"
	"os"

	"github.com/burubur/go-microservice/cmd/servd"
	"github.com/burubur/go-microservice/internal/platform/configuration"
)

const (
	inquiry     = "inquiry"
	payment     = "payment"
	checkStatus = "checkstatus"
	defaultPath = "configuration/app.yaml"
)

// Command is a main command which will instruct all of the availeble service
type Command struct{}

// New will instantiate a new instance of Command itself
func New() *Command {
	return &Command{}
}

// Execute will validate CLI args and then call Run method if it's a valid command
func (c *Command) Execute() {
	args := os.Args[1:]
	switch argsLen := len(args); {
	case argsLen == 1:
		c.Run(args)
	default:
		log.Println("our service currently handle 1 command only")
	}
}

// Run will serve the available service based on executed CLI command
func (c *Command) Run(args []string) {
	service := servd.New()
	c.loadAndValidateConfig()

	switch args[0] {
	case inquiry:
		service.Inquiry()
	case payment:
		service.Payment()
	case checkStatus:
		service.CheckStatus()
	default:
		log.Println("please specify the available command (inquiry, payment, checkstatus)")
	}
}

func (c *Command) loadAndValidateConfig() {
	err := configuration.Load(defaultPath)
	if err != nil {
		log.Fatalf("\nfailed to load configuration file, got: %v\n", err)
	}
	err = configuration.Validate()
	if err != nil {
		log.Fatalf("\ninvalid configuration values, got: %+v\n", err)
	}
}
