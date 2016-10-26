package main

import (
	"fmt"
	"os"

	"io/ioutil"
	"path"
	"sort"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type FileInfos []os.FileInfo
type ByName struct{ FileInfos }

func (fi ByName) Len() int {
	return len(fi.FileInfos)
}

func (fi ByName) Swap(i, j int) {
	fi.FileInfos[i], fi.FileInfos[j] = fi.FileInfos[j], fi.FileInfos[i]
}

func (fi ByName) Less(i, j int) bool {
	return fi.FileInfos[j].ModTime().Unix() < fi.FileInfos[i].ModTime().Unix()
}

func IsDirectory(name string) (isDir bool, err error) {
	fInfo, err := os.Stat(name) // FileInfo型が返る
	if err != nil {
		return false, err
	}
	return fInfo.IsDir(), nil
}

func main() {
	app := cli.NewApp()
	app.Name = "delete files"
	app.Usage = "Delete files!"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Value: "",
			Usage: "set target path.",
		},
		cli.BoolFlag{
			Name:  "delete",
			Usage: "whether delele file or not. (default is false)",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println(c.String("path"))
		fmt.Println(c.Bool("delete"))

		var argPath = c.String("path")
		var argDelete = c.Bool("delete")

		var curDir, _ = os.Getwd()
		curDir += "/"

		if argPath == "" {
			argPath = curDir
		}

		var dirName, filePattern = path.Split(argPath)

		if dirName == "" {
			dirName = curDir
		}

		var isDir, _ = IsDirectory(dirName + filePattern)

		if isDir == true {
			dirName = dirName + filePattern
			// filePattern = ""
			filePattern = "backuplog_*"
		}

		fileInfos, err := ioutil.ReadDir(dirName)

		if err != nil {
			// fmt.Errorf("Directory cannot read %s\n", err)
			// os.Exit(1)
			return errors.Wrap(err, "Directory cannot read") // Directory cannot read: open hoge: no such file or directory
		}

		sort.Sort(ByName{fileInfos})
		for _, fileInfo := range fileInfos {
			var findName = (fileInfo).Name()
			var matched = true

			if filePattern != "" {
				matched, _ = path.Match(filePattern, findName)
			}

			if matched == true {
				fmt.Printf("delete file is %s\n", findName)

				if argDelete == true {
					fmt.Println("delete!!")

					if err := os.Remove(argPath + "/" + findName); err != nil {
						return errors.Wrap(err, "Can not delete file")
					}
				}

			}
		}

		return nil
	}

	app.Run(os.Args)

}
