package main

import (
	"gofr.dev/pkg/gofr"

	depClient "zop.dev/client/deploy"
	depHndler "zop.dev/handler/deploy"
	deploySvc "zop.dev/service/deploy"
	dockerSvc "zop.dev/service/docker"
)

func main() {
	app := gofr.NewCMD()

	app.AddHTTPService("deployment-service", "https://api.stage.kops.dev")

	dClient := depClient.New()
	docker := dockerSvc.New()
	depSvc := deploySvc.New(docker, dClient)
	depHandler := depHndler.New(depSvc)

	app.SubCommand("version", func(_ *gofr.Context) (interface{}, error) {
		return "zop.dev cli version " + version, nil
	}, gofr.AddDescription("displays the installed zop.dev cli version"))

	app.SubCommand("deploy", depHandler.Deploy,
		gofr.AddDescription("builds and deploy code using a single command"))

	app.Run()
}
