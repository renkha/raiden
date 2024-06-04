package main

import (
	"raiden/internal/bootstrap"

	"github.com/sev-2/raiden"
	"github.com/sev-2/raiden/pkg/cli/generate"
	"github.com/sev-2/raiden/pkg/logger"
	"github.com/sev-2/raiden/pkg/resource"
	"github.com/sev-2/raiden/pkg/utils"
	"github.com/spf13/cobra"
)

func main() {
	f := resource.Flags{}

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			f.CheckAndActivateDebug(cmd)
			// load configuration
			if f.ProjectPath == "" {
				curDir, err := utils.GetCurrentDirectory()
				if err != nil {
					logger.Error(err)
					return
				}
				f.ProjectPath = curDir
			}

			config, err := raiden.LoadConfig(nil)
			if err != nil {
				logger.Error(err)
				return
			}

			// register app resource
			bootstrap.RegisterRpc()
			bootstrap.RegisterRoles()
			bootstrap.RegisterModels()
			bootstrap.RegisterStorages()

			if err := resource.Import(&f, config); err != nil {
				logger.Error(err)
				return
			}

			if err = generate.Run(&f.Generate, config, f.ProjectPath, false); err != nil {
				logger.Error(err)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&f.Verbose, "verbose", "v", false, "verbose mode")
	cmd.Flags().StringVarP(&f.ProjectPath, "project-path", "p", "", "set project path")
	cmd.Flags().BoolVarP(&f.RpcOnly, "rpc-only", "", false, "import rpc only")
	cmd.Flags().BoolVarP(&f.RolesOnly, "roles-only", "r", false, "import roles only")
	cmd.Flags().BoolVarP(&f.ModelsOnly, "models-only", "m", false, "import models only")
	cmd.Flags().StringVarP(&f.AllowedSchema, "schema", "s", "", "set allowed schema to import, use coma separator for multiple schema")

	f.Generate.Bind(cmd)

	cmd.Execute()
}
