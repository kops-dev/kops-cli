package main

import (
	"gofr.dev/pkg/gofr"

	depClient "kops.dev/client/deploy"
	depHndler "kops.dev/handler/deploy"
	deploySvc "kops.dev/service/deploy"
	dockerSvc "kops.dev/service/docker"
)

func main() {
	app := gofr.NewCMD()

	app.AddHTTPService("deployment-service", "https://api.kops.dev")

	dClient := depClient.New()
	docker := dockerSvc.New()
	depSvc := deploySvc.New(docker, dClient)
	depHandler := depHndler.New(depSvc)

	app.SubCommand("version", func(_ *gofr.Context) (interface{}, error) {
		return "kops cli version " + version, nil
	}, gofr.AddDescription("displays the installed kops version"))

	app.SubCommand("deploy", depHandler.Deploy,
		gofr.AddDescription("builds and deploy code using a single command"))

	app.Run()
}
