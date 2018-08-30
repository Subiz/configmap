package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "configmap"
	app.Usage = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: "",
			Usage: "vault address",
		},
		cli.StringFlag{
			Name:  "token",
			Value: "",
			Usage: "vault token",
		},
	}

	app.Action = run
	l := log.New(os.Stderr, "", 0)
	if err := app.Run(os.Args); err != nil {
		l.Fatal(err)
	}
}

func loadConfigMap(name string) ([]Config, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	obj := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(data), &obj); err != nil {
		return nil, err
	}

	return ParseConfigMap(obj, os.Environ()), nil
}

func run(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing file name")
	}
	name := c.Args().Get(0)

	// parse parameters
	addr := strings.TrimSpace(c.String("addr"))
	token := strings.TrimSpace(c.String("token"))
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
		if addr == "" {
			addr = "http://localhost:8200"
		}
	}

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

	out, err := parse(configs, vaultdata)
	fmt.Println(out)

	if err != nil {
		return err
	}
	return nil
}
