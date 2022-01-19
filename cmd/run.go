/*
Copyright © 2021 Red Hat D&O Tools Development Team

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

// Package cmd is the package for all command line related things
package cmd

import (
	"github.com/RedHatQE/tfacon/common"
	"github.com/RedHatQE/tfacon/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the info retrival and get the prediction from TFA",
	Long:  `run the info retrival and get the prediction from TFA`,
	Run: func(cmd *cobra.Command, args []string) {
		common.PrintHeader(rootCmd.Version)
		core.Run(viperRun, viperConfig)
	},
}
var viperRun *viper.Viper

func init() {
	rootCmd.AddCommand(runCmd)

	viperRun = viper.New()
	initConfig(viperRun, runCmd, cmdInfoList)
}
