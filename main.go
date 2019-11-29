package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gomove"
	app.Usage = "Move Golang packages to a new path."
	app.Version = "0.2.17"
	app.ArgsUsage = "[old path] [new path]"
	app.Authors = append(app.Authors, &cli.Author{Name: "Kaushal Subedi", Email: "<kaushal@subedi.co>"})

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "dir, d",
			Value: "./",
			Usage: "directory to scan. files under vendor/ are ignored",
		},
		&cli.StringFlag{
			Name:  "file, f",
			Usage: "only move imports in a file",
		},
		&cli.StringFlag{
			Name:  "safe-mode, s",
			Value: "false",
			Usage: "run program in safe mode (comments will be wiped)",
		},
		// The reason why this flag applies only to `safe-mode` is that in non-`safe-mode` we do
		// substring replaces anyway.
		&cli.StringFlag{
			Name:  "prefix, p",
			Value: "false",
			Usage: "interpret 'from' and 'to' arguments as import path prefixes rather than the entire paths (applies only to 'safe-mode')",
		},
	}

	app.Action = func(c *cli.Context) error {
		file := c.String("file")
		dir := c.String("dir")
		from := c.Args().Get(0)
		to := c.Args().Get(1)

		if file != "" {
			ProcessFile(file, from, to, c)
		} else {
			ScanDir(dir, from, to, c)
		}

		return nil
	}

	app.Run(os.Args)
}

// ScanDir scans a directory for go files and
func ScanDir(dir string, from string, to string, c *cli.Context) {
	// If from and to are not empty scan all files
	if from != "" && to != "" {
		// construct the path of the vendor dir of the project for prefix matching
		vendorDir := path.Join(dir, "vendor")
		// Scan directory for files
		filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
			// ignore vendor path
			if matched := strings.HasPrefix(filePath, vendorDir); matched {
				return nil
			}
			// Only process go files
			if path.Ext(filePath) == ".go" {
				ProcessFile(filePath, from, to, c)
			}

			return nil
		})

	} else {
		cli.ShowAppHelp(c)
	}

}

// ProcessFile processes file according to what mode is chosen
func ProcessFile(filePath string, from string, to string, c *cli.Context) {
	if c.String("safe-mode") == "true" {
		ProcessFileAST(filePath, from, to, c)
	} else {
		ProcessFileNative(filePath, from, to)
	}
}
