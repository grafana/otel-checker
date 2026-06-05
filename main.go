package main

import (
	"log"

	"github.com/grafana/otel-checker/checks"
	"github.com/grafana/otel-checker/checks/utils"
	"github.com/grafana/otel-checker/checks/webserver"
)

func main() {
	commands := utils.GetArguments()
	messages := checks.RunAllChecks(commands)

	if !commands.WebServer {
		return
	}

	if err := webserver.Run(":8080", messages); err != nil {
		log.Fatal(err)
	}
}
