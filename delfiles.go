package main

import (
	"fmt"
	"os"
	"time"

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

func isOld(fileTimeStamp time.Time, elaspedDays int) bool {
	base := time.Now().AddDate(0, 0, elaspedDays*-1)
	if base.Before(fileTimeStamp) {
		return false // file is new.
	}

	return true // file is old.
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
		cli.IntFlag{
			Name:  "days",
			Usage: "set elasped days of delete files.",
		},
	}

	app.Action = func(c *cli.Context) error {
		var argPath = c.String("path")
		var argDelete = c.Bool("delete")
		var argDays = c.Int("days")

		var curDir, _ = os.Getwd()
		curDir += "/"

		if argPath == "" {
			argPath = curDir
		}

		if argDays == 0 {
			argDays = 7
		}

		if argDelete == false {
			fmt.Println("Dry-Run")
		}

		var dirName, filePattern = path.Split(argPath)

		if dirName == "" {
			dirName = curDir
		}

		var isDir, _ = IsDirectory(dirName + filePattern)

		if isDir == true {
			dirName = dirName + filePattern
			filePattern = "backuplog_*"
		}

		fileInfos, err := ioutil.ReadDir(dirName)

		if err != nil {
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

				// check timestamp
				if isOld(fileInfo.ModTime(), argDays) == false {
					continue
				}

				fmt.Printf("delete %s timestamp: %s\n", findName, fileInfo.ModTime())

				if argDelete == true {
					if err := os.Remove(argPath + "/" + findName); err != nil {
						return errors.Wrap(err, "Can not delete file")
					}

					fmt.Println("deleted!!")
				}

			}
		}

		return nil
	}

	app.Run(os.Args)

}
