package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grafana/otel-checker/checks"
	"github.com/grafana/otel-checker/checks/output"
	"github.com/grafana/otel-checker/checks/utils"
	"github.com/grafana/otel-checker/checks/webserver"
)

func main() {
	commands := utils.GetArguments()
	reporter := checks.Run(commands)

	if err := output.Render(os.Stdout, reporter, commands.Format); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !commands.WebServer {
		return
	}

	if err := webserver.Run(":8080", reporter.Results()); err != nil {
		log.Fatal(err)
	}
}
