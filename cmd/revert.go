/*
Copyright © 2023 Red Hat D&O Tools Development Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/RedHatQE/tfacon/common"
	"github.com/RedHatQE/tfacon/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// revertCmd represents the revert command.
var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "revert will revert all test items from the input launch id and project to the original status",
	Long:  `revert will revert all test items from the input launch id and project to the original status`,
	Run: func(cmd *cobra.Command, args []string) {
		common.PrintHeader(rootCmd.Version)
		core.Revert(viperRevert, viperConfig)
	},
}
var viperRevert *viper.Viper

func init() {
	rootCmd.AddCommand(revertCmd)

	viperRevert = viper.New()
	initConfig(viperRevert, revertCmd, cmdInfoList)
}
