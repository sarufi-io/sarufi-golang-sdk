package main

import (
	"fmt"
	"os"

	"github.com/sarufi-io/sarufi-golang-sdk/cli"
)

func main() {
	manager := cli.NewManager(&cli.Credentials{
		Username: "username",
		Password: "password",
	})

	if err := manager.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
