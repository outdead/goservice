package main

import "github.com/outdead/goservice/internal/app"

// ServiceName contains the name of the service. Displayed in logs and when help
// command is called.
const ServiceName = "goservice"

// ServiceVersion contains the service version number in the semantic versioning
// format (http://semver.org/). Displayed in logs and when help command is
// called. During service compilation, you can pass the version value using the
// flag `-ldflags "-X main.Version=${VERSION}"`.
var ServiceVersion = "0.0.0-develop"

func main() {
	app.New(ServiceName, ServiceVersion).Run()
}
