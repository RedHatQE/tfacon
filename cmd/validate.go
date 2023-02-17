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
	"fmt"

	"github.com/RedHatQE/tfacon/common"
	"github.com/RedHatQE/tfacon/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validateCmd represents the validate command.
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate if the parameter is valid and if the urls are accessible",
	Long:  `validate if the parameter is valid and if the urls are accessible`,
	Run: func(cmd *cobra.Command, args []string) {
		con := core.GetInfo(viperValidate)
		fmt.Println(con.String())
		success, err := core.Validate(con, viperConfig)
		if err == nil && success {
			common.PrintGreen("Validation Passed!")
		} else {
			fmt.Println()
			common.PrintRed(fmt.Sprintf("There is an error during validation: \n%s", err))
		}
	},
}

var viperValidate *viper.Viper

func init() {
	rootCmd.AddCommand(validateCmd)

	viperValidate = viper.New()
	initConfig(viperValidate, validateCmd, cmdInfoList)
}
