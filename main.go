package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "configmap"
	app.Usage = "configmap"
	app.Version = "1.0.11"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "take input as config file",
		},
		cli.StringFlag{
			Name:  "format",
			Value: "docker",
			Usage: "output format, can be bash, docker",
		},
	}
	app.Action = run
	l := log.New(os.Stderr, "", 0)
	if err := app.Run(os.Args); err != nil {
		l.Fatal(err)
	}
}

func loadConfigMap(name string) ([]Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	obj := make(map[interface{}]interface{})
	dec := yaml.NewDecoder(f)
	for {
		err := dec.Decode(&obj)
		if err == nil {
			continue
		}

		if err == io.EOF {
			break
		}

		return nil, err
	}

	return ParseConfigMap(obj, os.Environ()), nil
}

func run(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.ShowAppHelp(c)
	}
	name := c.Args().Get(0)

	format := strings.TrimSpace(c.String("format"))

	configs, err := loadConfigMap(name)
	if err != nil {
		return err
	}

	paths, fields := make([]string, 0), make([]string, 0)
	for _, c := range configs {
		paths = append(paths, c.ConfigPath)
		fields = append(fields, c.ConfigField)
	}

	configpath := strings.TrimSpace(c.String("config-file"))
	data, err := readFile(configpath, paths, fields)

	if err != nil {
		return err
	}

	out, err := parse(configs, data, format)
	fmt.Println(out)
	return err
}
