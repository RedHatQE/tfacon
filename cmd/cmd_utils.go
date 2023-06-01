// Package cmd is the package for all command line related things
package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/RedHatQE/tfacon/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdInfoList []map[string]string = []map[string]string{
	{
		"cmdName":        "tfa-url",
		"valName":        "TFA_URL",
		"defaultVal":     "default value for tfa url",
		"cmdDescription": "The url to the TFA Classifier",
	},
	{
		"cmdName":        "re-url",
		"valName":        "RE_URL",
		"defaultVal":     "default value for Recommendation Engine url",
		"cmdDescription": "The url to the Recommendation Engine",
	},
	{
		"cmdName":        "platform-url",
		"valName":        "PLATFORM_URL",
		"defaultVal":     "default value for platform url",
		"cmdDescription": "The url to the test platform (example: https://reportportal-<your_domain>.com)",
	},
	{
		"cmdName":        "connector-type",
		"valName":        "CONNECTOR_TYPE",
		"defaultVal":     "RPCon",
		"cmdDescription": "The type of connector you want to use (example: RPCon, PolarionCon, JiraCon)",
	},
	{
		"cmdName":        "launch-id",
		"valName":        "LAUNCH_ID",
		"defaultVal":     "",
		"cmdDescription": "The launch id of report portal",
	},
	{
		"cmdName":        "launch-uuid",
		"valName":        "LAUNCH_UUID",
		"defaultVal":     "",
		"cmdDescription": "The launch uuid of report portal",
	},
	{
		"cmdName":        "project-name",
		"valName":        "PROJECT_NAME",
		"defaultVal":     "",
		"cmdDescription": "The project name of report portal",
	},
	{
		"cmdName":        "team-name",
		"valName":        "TEAM_NAME",
		"defaultVal":     "",
		"cmdDescription": "Your team name",
	},
	{
		"cmdName":        "auth-token",
		"valName":        "AUTH_TOKEN",
		"defaultVal":     "",
		"cmdDescription": "The AUTH_TOKEN of report portal",
	},
	{
		"cmdName":        "launch-name",
		"valName":        "LAUNCH_NAME",
		"defaultVal":     "",
		"cmdDescription": "The launch name of the launch in report portal",
	},
}

func initConfig(viper *viper.Viper, cmd *cobra.Command, cmdInfoList []map[string]string) {
	switch {
	case os.Getenv("TFACON_YAML_PATH") != "":
		index := strings.LastIndex(os.Getenv("TFACON_YAML_PATH"), "/")
		path := os.Getenv("TFACON_YAML_PATH")[:index]
		viper.AddConfigPath(path)

		configName := strings.Split(os.Getenv("TFACON_YAML_PATH")[index+1:], ".")
		viper.SetConfigName(configName[0])

	case common.FileExist("./tfacon.yml") || common.FileExist("./tfacon.yaml"):
		viper.AddConfigPath(".")
		viper.SetConfigName("tfacon")

	default:

		_, err := os.Create("tfacon.yml")
		common.HandleError(err, "nopanic")
		viper.AddConfigPath(".")
		viper.SetConfigName("tfacon")
	}

	viper.AutomaticEnv()

	for _, v := range cmdInfoList {
		initViperVal(cmd, viper, v["cmdName"], v["valName"], v["defaultVal"], v["cmdDescription"])
	}

	err := viper.ReadInConfig()

	common.HandleError(err, "nopanic")
}

func initViperVal(cmd *cobra.Command, viper *viper.Viper, cmdName, valName, defaultVal, cmdDescription string) {
	if viper.GetString(valName) == "" {
		viper.SetDefault(valName, defaultVal)
	} else {
		viper.SetDefault(valName, viper.GetString(valName))
	}

	cmd.PersistentFlags().StringP(cmdName, "", viper.GetString(valName), cmdDescription)
	err := viper.BindPFlag(valName, cmd.PersistentFlags().Lookup(cmdName))
	common.HandleError(err, "nopanic")
}

func InitTFAConfigFile(viper *viper.Viper) {
	var file []byte

	var err error

	if os.Getenv("TFACON_CONFIG_PATH") != "" {
		file, err = ioutil.ReadFile(os.Getenv("TFACON_CONFIG_PATH"))
	} else if common.FileExist("./tfacon.cfg") {
		file, err = ioutil.ReadFile("./tfacon.cfg")
	}

	defer func() {
		if r := recover(); r != nil {
			_, err = os.Create("tfacon.cfg")
			common.HandleError(err, "nopanic")
		}
	}()
	common.HandleError(err, "nopanic")
	viper.SetConfigType("ini")
	viper.SetDefault("config.concurrency", true)
	viper.SetDefault("config.retry_times", 20)
	viper.SetDefault("config.add_attributes", false)
	err = viper.ReadConfig(bytes.NewBuffer(file))
	common.HandleError(err, "nopanic")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "You can add this tag to print more detailed info")
	err = viper.BindPFlag("config.verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolP("re", "r", false,
		"You can add this tag to let tfacon add Recommendation Engine result to comment section")
	err = viper.BindPFlag("config.re", rootCmd.PersistentFlags().Lookup("re"))
	common.HandleError(err, "nopanic")
}

func initWorkspace() {
	pwd, err := os.Getwd()
	common.HandleError(err, "nopanic")
	fmt.Println(pwd)

	err = os.Mkdir("tfacon_workspace", fs.ModePerm)

	common.HandleError(err, "nopanic")
	err = os.Mkdir("/tmp/.tfacon", fs.ModePerm)
	common.HandleError(err, "nopanic")
	err = os.Chdir("/tmp/.tfacon/")
	common.HandleError(err, "nopanic")

	cmd := exec.Command("git", "clone", "https://github.com/RedHatQE/tfacon.git")

	err = cmd.Run()
	common.HandleError(err, "nopanic")
	err = os.Chdir("/tmp/.tfacon/tfacon")

	common.HandleError(err, "nopanic")

	cmd = exec.Command("mv", "examples", pwd+"/tfacon_workspace")

	err = cmd.Run()
	common.HandleError(err, "nopanic")

	cmd = exec.Command("rm", "/tmp/.tfacon", "-rf")

	err = cmd.Run()
	common.HandleError(err, "nopanic")
}
