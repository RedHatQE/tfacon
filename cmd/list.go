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

// Package cmd is the package for all command line related things
package cmd

import (
	"log"

	"github.com/RedHatQE/tfacon/common"
	"github.com/RedHatQE/tfacon/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all constructed information",
	Long:  `list all information constructed from tfacon.yml/environment variables/cli`,
	Run: func(cmd *cobra.Command, args []string) {
		common.PrintHeader(rootCmd.Version)
		log.Println("Printing the constructed information")
		con := core.GetInfo(viperList)
		common.PrintGreen(con.String())
	},
}
var viperList *viper.Viper

func init() {
	rootCmd.AddCommand(listCmd)

	viperList = viper.New()
	initConfig(viperList, listCmd, cmdInfoList)
}
