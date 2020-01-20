package cmds

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mdevilliers/vault-penknife/pkg/vault"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func RegisterCopyCmd(client *api.Client) func() (*cobra.Command, error) {

	return func() (*cobra.Command, error) {

		var source, destination string
		var dryrun bool

		copyCommand := &cobra.Command{
			Use:   "copy",
			Short: "Copy a tree of Vault secrets into another.",
			RunE: func(cmd *cobra.Command, args []string) error {

				allKeys, err := vault.Walk(client, source)
				if err != nil {
					return errors.Wrap(err, "error collecting source keys")
				}

				for _, k := range allKeys {

					newK := strings.ReplaceAll(k, source, destination)

					if dryrun {
						fmt.Println("[dry-run]", k, "copy to", newK)
					} else {
						from, err := client.Logical().Read(k)

						if err != nil {
							return errors.Wrapf(err, "error reading value of source key: %s", k)
						}

						_, err = client.Logical().Write(newK, from.Data)

						if err != nil {
							return errors.Wrapf(err, "error writing value to dest key: %s", newK)
						}
					}
				}
				return nil
			},
		}

		copyCommand.Flags().StringVar(&source, "source", source, "source path e.g. secret/foo")
		copyCommand.Flags().StringVar(&destination, "dest", destination, "destination path e.g. secret/bar")
		copyCommand.Flags().BoolVar(&dryrun, "dry-run", dryrun, "show me what you would have done - no effects")

		err := markFlagsRequired(copyCommand, "source", "dest")

		if err != nil {
			return nil, err
		}
		return copyCommand, nil
	}
}

func markFlagsRequired(cmd *cobra.Command, flags ...string) error {

	for _, f := range flags {
		err := cmd.MarkFlagRequired(f)

		if err != nil {
			return err
		}

	}
	return nil
}
