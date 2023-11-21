package core

import (
	"fmt"
	"net/http"

	"github.com/RedHatQE/tfacon/common"
	"github.com/RedHatQE/tfacon/connectors"
	"github.com/spf13/viper"
)

// TFACon is the general interface for all TFA Classifer
// connectors, any connector to any test management platform
// should inpement this interface.
type TFACon interface {
	GetAllTestInfos() map[string]string
	GetAllTestIds() []string
	BuildUpdatedList(ids []string, concurrent bool, add_attributes bool, re bool, auto_finalize_defect_type bool, auto_finalization_thredshold float32) common.GeneralUpdatedList
	UpdateAll(common.GeneralUpdatedList, bool)
	String() string
	InitConnector()
	Validate(verbose bool) (bool, error)
	RevertUpdatedList(verbose bool) common.GeneralUpdatedList
}

func Revert(viperRevert, viperConfig *viper.Viper) {
	var con TFACon = GetCon(viperRevert)

	runHelper(viperConfig, con.GetAllTestIds(), con, "revert")
}

// Run method is the run operation for any type of connector that
// implements TFACon interface.
func Run(viperRun, viperConfig *viper.Viper) {
	var con TFACon = GetCon(viperRun)

	con.InitConnector()
	ids := con.GetAllTestIds()
	retry := viperConfig.GetInt("config.retry_times")
	original_retry := retry
	for len(ids) != 0 && retry > 0 {
		runHelper(viperConfig, ids, con, "run")
		ids = con.GetAllTestIds()
		fmt.Printf("This is the %d retry\n", original_retry-retry+1)
		retry--
	}
	if len(con.GetAllTestIds()) != 0 {
		panic("Update Failed, check the tfa-c/tfa-r service")
	}
}

func runHelper(viperConfig *viper.Viper, ids []string, con TFACon, operation string) {
	if len(ids) == 0 {
		return
	}

	var updated_list_of_issues common.GeneralUpdatedList

	switch operation {
	case "run":
		updated_list_of_issues = con.BuildUpdatedList(ids, viperConfig.GetBool("config.concurrency"),
			viperConfig.GetBool("config.add_attributes"), viperConfig.GetBool("config.re"), viperConfig.GetBool("config.auto_finalize_defect_type"), float32(viperConfig.GetFloat64("config.auto_finalization_thredshold")))
	case "revert":
		updated_list_of_issues = con.RevertUpdatedList(viperConfig.GetBool("config.verbose"))
	default:
		updated_list_of_issues = con.RevertUpdatedList(viperConfig.GetBool("config.verbose"))
	}
	// Doing this because the api can only take 20 items per request
	con.UpdateAll(updated_list_of_issues, viperConfig.GetBool("config.verbose"))
}

// GetInfo method is the get info operation for any type of connector that
// implements TFACon interface.
func GetInfo(viper *viper.Viper) TFACon {
	con := GetCon(viper)

	return con
}

// Validate method is the validate operation for any type of connector that
// implements TFACon interface.
func Validate(con TFACon, viper *viper.Viper) (bool, error) {
	success, err := con.(*connectors.RPConnector).Validate(viper.GetBool("config.verbose"))
	if err != nil {
		err = fmt.Errorf("validate error: %w", err)
	}

	return success, err
}

// GetCon method is the get con operation for any type of connector that
// implements TFACon interface, it returns the actual tfa connector instance.
func GetCon(viper *viper.Viper) TFACon {
	var con TFACon

	switch viper.Get("CONNECTOR_TYPE") {
	case "RPCon":
		con = &connectors.RPConnector{Client: &http.Client{}}
		err := viper.Unmarshal(con)

		common.HandleError(err, "panic")
	// case "POLCon":
	// 	con = RPConnector{}
	// case "JiraCon":
	// 	con = RPConnector{}
	default:
		con = &connectors.RPConnector{Client: &http.Client{}}

		err := viper.Unmarshal(con)

		common.HandleError(err, "panic")
	}

	return con
}
