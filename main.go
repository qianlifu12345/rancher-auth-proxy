package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-auth-filter-service/manager"
	"github.com/rancher/rancher-auth-filter-service/service"
	"github.com/urfave/cli"
)

//VERSION for Rancher Authantication Filter Service
var VERSION = "v0.1.0-dev"

func main() {

	///init parsing command line
	app := cli.NewApp()
	app.Name = "rancher-auth-filter-service"
	app.Version = "v0.1.0-dev"
	app.Usage = "Rancher Authantication Filter Service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "rancherUrl",
			Value:  "http://54.169.69.238:8080/",
			Usage:  "Rancher server url",
			EnvVar: "RANCHER_SERVER_URL",
		},
		cli.StringFlag{
			Name:   "localport",
			Value:  "8092",
			Usage:  "Local server port ",
			EnvVar: "LOCAL_VALIDATION_FILTER_PORT",
		},
	}

	app.Action = func(c *cli.Context) error {

		manager.URL = c.String("rancherUrl")
		manager.Port = c.String("localport")

		logrus.Infof("Starting token validation service")
		logrus.Infof("Rancher server URL:" + manager.URL + ". The validation filter server running on local port:" + manager.Port)
		//create mux router
		router := service.NewRouter()
		http.Handle("/", router)
		serverString := ":" + manager.Port
		//start local server
		logrus.Fatal(http.ListenAndServe(serverString, nil))
		return nil
	}

	app.Run(os.Args)

}
