package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "delete files"
	app.Usage = "Delete files!"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Value: "hoge",
			Usage: "aaaaaaaa",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("run")

		return nil
	}

	app.Run(os.Args)

}
