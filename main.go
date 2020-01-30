package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/kelseyhightower/envconfig"
	"github.com/mdevilliers/go/cli"
	"github.com/mdevilliers/vault-penknife/internal/cmds"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "vault-penknife",
}

type config struct {
	Token string `envconfig:"VAULT_PENKNIFE_TOKEN" default:""`
}

var (
	globalConfig = config{}
)

func initConfig() {
	err := envconfig.Process("", &globalConfig)

	if err != nil {
		panic(err)
	}
}

func main() {
	initConfig()

	cfg := api.DefaultConfig()
	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	client.SetToken(globalConfig.Token)

	err = cli.RegisterCommands(rootCmd,
		cmds.RegisterCopyCmd(client))

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
