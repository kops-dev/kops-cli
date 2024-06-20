package main

import "gofr.dev/pkg/gofr"

func main() {
	app := gofr.NewCMD()

	app.SubCommand("version", func(ctx *gofr.Context) (interface{}, error) {
		return "kops cli version " + version, nil
	}, gofr.AddDescription("displays the installed kops version"))

	app.Run()
}
