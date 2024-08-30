package main

import (
	"gofr.dev/pkg/gofr"
	gofrSvc "gofr.dev/pkg/gofr/service"

	"kops.dev/deployment-service/client/kops"
	"kops.dev/deployment-service/handler/deploy"
	depSvc "kops.dev/deployment-service/service/deploy"
	"kops.dev/deployment-service/service/upload"
)

func main() {
	app := gofr.New()

	kopsClient := gofrSvc.NewHTTPService(app.Config.Get("KOPS_SERVICE_ADDR"), app.Logger(), app.Metrics(),
		&gofrSvc.APIKeyConfig{APIKey: app.Config.Get("KOPS_API_KEY")})

	kopsSvc := kops.New(kopsClient)
	uploadSvc := upload.New(kopsSvc)
	deploySvc := depSvc.New(kopsSvc)

	deployHndlr := deploy.New(uploadSvc, deploySvc)

	app.POST("/deploy", deployHndlr.UploadImage)

	app.Run()
}
