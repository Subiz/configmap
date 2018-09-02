package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "configmap"
	app.Usage = "configmap"
	app.Version = "1.0.10"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "format",
			Value: "docker",
			Usage: "output format, can be bash, docker",
		},
		cli.StringFlag{
			Name:  "addr",
			Value: os.Getenv("VAULT_ADDR"),
			Usage: "vault address (VAULT_ADDR)",
		},
		cli.StringFlag{
			Name:  "token",
			Value: os.Getenv("VAULT_TOKEN"),
			Usage: "vault token (VAULT_TOKEN)",
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
	for dec.Decode(&obj) == nil {
	}

	return ParseConfigMap(obj, os.Environ()), nil
}

func run(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.ShowAppHelp(c)
	}
	name := c.Args().Get(0)

	format := strings.TrimSpace(c.String("format"))
	addr := strings.TrimSpace(c.String("addr"))
	token := strings.TrimSpace(c.String("token"))

	configs, err := loadConfigMap(name)
	if err != nil {
		return err
	}

	paths, fields := make([]string, 0), make([]string, 0)
	for _, c := range configs {
		paths = append(paths, c.VaultPath)
		fields = append(fields, c.VaultField)
	}

	vaultdata, err := readVault(addr, token, paths, fields)
	if err != nil {
		return err
	}

	out, err := parse(configs, vaultdata, format)
	fmt.Println(out)
	return err
}
