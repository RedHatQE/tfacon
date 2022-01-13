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
	"fmt"
	"os"

	"github.com/RedHatQE/tfacon/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "tfacon",
	Short: "A connector tool to connect testing platform and TFA Classifier",
	Long:  `A connector tool to connect testing platform and TFA Classifier`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			common.HandleError(err)
			os.Exit(0)
		}
	},
	Version: "1.0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	viperConfig *viper.Viper
	cfg         map[string]map[string]string
)

func init() {
	viperConfig = viper.New()
	initTFAConfigFile(viperConfig)
	err := viperConfig.Unmarshal(&cfg)
	common.HandleError(err)
}
