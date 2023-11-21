// Package connectors is the package for all connectors struct and its
// coressponding methods
package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/RedHatQE/tfacon/common"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Issues is a list of IssueItem.
type Issues []IssueItem
type REResult struct {
	LaunchID     string `json:"launch_id"`
	LaunchName   string `json:"launch_name"`
	TestItemId   string `json:"test_item_id"`
	TestItemName string `json:"test_name"`
	IssueType    string `json:"issue_name"`
	TicketNumber string `json:"ticket_number"`
	TicketURL    string `json:"ticket_url"`
	Team         string `json:"TEAM"`
}

// IssueItem is a smallest unit
// in the reuquest body for describing
// the issue to be updated.
type IssueItem struct {
	Issue      IssueInfo `json:"issue"`
	TestItemID string    `json:"testItemId"`
}

// IssueInfo is the data structure describing
// the issue information for the request body.
type IssueInfo struct {
	IssueType            string        `json:"issueType"`
	Comment              string        `json:"comment"`
	AutoAnalyzed         bool          `json:"autoAnalyzed"`
	IgnoreAnalyzer       bool          `json:"ignoreAnalyzer"`
	ExternalSystemIssues []interface{} `json:"externalSystemIssues"`
}

// UpdatedList is the finished final list
// of the request body that contains all finalized
// information about issues(test items).
type UpdatedList struct {
	IssuesList Issues `json:"issues"`
}

// GetSelf method returns an GenralUpdatedlist.
func (u UpdatedList) GetSelf() common.GeneralUpdatedList {
	return u
}

// RPConnector is the class for describing
// the RPConnector engine.
type RPConnector struct {
	LaunchID    string `mapstructure:"LAUNCH_ID" json:"launch_id"`
	LaunchUUID  string `mapstructure:"LAUNCH_UUID" json:"uuid"`
	LaunchName  string `mapstructure:"LAUNCH_NAME" json:"launch_name"`
	ProjectName string `mapstructure:"PROJECT_NAME" json:"project_name"`
	TeamName    string `mapstructure:"TEAM_NAME" json:"team_name"`
	AuthToken   string `mapstructure:"AUTH_TOKEN" json:"auth_token"`
	RPURL       string `mapstructure:"PLATFORM_URL" json:"platform_url"`
	Client      *http.Client
	TFAURL      string `mapstructure:"TFA_URL" json:"tfa_url"`
	REURL       string `mapstructure:"RE_URL" json:"re_url"`
}

// Validate method validates against the input from
// yaml file, cli flag, and environment variable.
func (c *RPConnector) Validate(verbose bool) (bool, error) {
	fmt.Print("Validating....\n")

	validateRPURLAndAuthToken, err := c.validateRPURLAndAuthToken(verbose)
	if err != nil {
		return false, err
	}

	validateTFA, err := c.validateTFAURL(verbose)
	if err != nil {
		return false, err
	}

	projectnameNotEmpty, err := c.validateProjectName(verbose)
	if err != nil {
		return false, err
	}

	launchinfoNotEmpty, err := c.validateLaunchInfo(verbose)
	if err != nil {
		return false, err
	}

	ret := validateRPURLAndAuthToken && validateTFA && projectnameNotEmpty && launchinfoNotEmpty

	return ret, err
}

func (c *RPConnector) validateTFAURL(verbose bool) (bool, error) {
	body := `{"data": {"id": "123", "project": "rhv", "messages": ""}}`

	_, success, err := common.SendHTTPRequest(context.Background(), http.MethodPost,
		c.TFAURL, "", bytes.NewBuffer([]byte(body)), c.Client)
	if err != nil {
		err = fmt.Errorf("validate tfa url failed: %w", err)
		common.HandleError(err, "nopanic")
	}

	if verbose {
		fmt.Printf("TFAURLValidate: %t\n", success)
	}

	return success, err
}

func (c *RPConnector) validateLaunchInfo(verbose bool) (bool, error) {
	launchinfoNotEmpty := c.LaunchID != "" || c.LaunchName != "" || c.LaunchUUID != ""

	if verbose {
		fmt.Printf("lauchinfoValidate: %t\n", launchinfoNotEmpty)
	}

	if !launchinfoNotEmpty {
		err := errors.Errorf("%s", "You need to input launch id or launch name")

		return false, err
	}

	return true, nil
}

func (c *RPConnector) validateProjectName(verbose bool) (bool, error) {
	projectnameNotEmpty := c.ProjectName != ""
	if verbose {
		fmt.Printf("projectnameValidate: %t\n", projectnameNotEmpty)
	}

	if !projectnameNotEmpty {
		err := errors.Errorf("%s", "You need to input project name")

		return false, err
	}

	return true, nil
}

func (c *RPConnector) validateRPURLAndAuthToken(verbose bool) (bool, error) {
	_, success, err := common.SendHTTPRequest(context.Background(), http.MethodGet,
		c.RPURL+"/api/v1/project/list", c.AuthToken, bytes.NewBuffer(nil), c.Client)
	if err != nil {
		err = fmt.Errorf("validate rp url and auth token failed: %w", err)
		common.HandleError(err, "nopanic")
	}

	if verbose {
		fmt.Printf("RPURLValidate: %t\n", success)
	}

	return success, err
}

// String method is a to string method for the
// RPConnector instance.
func (c RPConnector) String() string {
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	str := ""

	for i := 0; i < v.NumField(); i++ {
		str += fmt.Sprintf("%s: \t %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}

	return str
}

// UpdateAll method is the interface method for tfacon interface
// it update the test items with the built updated_list.
func (c *RPConnector) UpdateAll(updatedListOfIssues common.GeneralUpdatedList, verbose bool) {
	if len(updatedListOfIssues.GetSelf().(UpdatedList).IssuesList) == 0 {
		return
	}

	json_updated_list_of_issues, _ := json.Marshal(updatedListOfIssues)

	log.Println("Updating All Test Items With Predictions...")

	url := fmt.Sprintf("%s/api/v1/%s/item", c.RPURL, c.ProjectName)
	method := "PUT"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(json_updated_list_of_issues)

	data, success, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	if err != nil {
		common.HandleError(errors.Errorf("Updated All failed: %s", err), "nopanic")
	}

	if verbose {
		fmt.Printf("This is the return info from update: %v\n", string(data))
	}

	if success {
		fmt.Println()
		common.PrintGreen("Updated All Test Items Successfully!")
	} else {
		common.PrintRed("Updated Failed!")
	}
}

// BuildUpdatedList method is a interface method for tfacon interface
// it builds a list of issues, it returns GeneralUpdatedList.
func (c *RPConnector) BuildUpdatedList(ids []string,
	concurrent bool, add_attributes bool, re bool, auto_finalize_defect_type bool, auto_finalization_thredshold float32) common.GeneralUpdatedList {
	return UpdatedList{IssuesList: c.BuildIssues(ids, concurrent, add_attributes, re, auto_finalize_defect_type, auto_finalization_thredshold)}
}

// BuildIssues method build the issue struct.
func (c *RPConnector) BuildIssues(ids []string, concurrent bool, add_attributes bool, re bool, auto_finalize_defect_type bool, auto_finalization_thredshold float32) Issues {
	issues := Issues{}

	if concurrent {
		return c.BuildIssuesConcurrent(ids, add_attributes, re, auto_finalize_defect_type, auto_finalization_thredshold)
	}

	for _, id := range ids {
		log.Printf("Getting prediction of test item(id): %s\n", id)
		issues = append(issues, c.BuildIssueItemHelper(id, add_attributes, re, auto_finalize_defect_type, auto_finalization_thredshold))
	}

	return issues
}

// BuildIssuesConcurrent methods builds the issues struct concurrently.
func (c *RPConnector) BuildIssuesConcurrent(ids []string, add_attributes bool, re bool, auto_finalize_defect_type bool, auto_finalization_thredshold float32) Issues {
	issues := Issues{}
	issuesChan := make(chan IssueItem, len(ids))
	idsChan := make(chan string, len(ids))
	exitChan := make(chan bool, len(ids))

	go func() {
		for _, id := range ids {
			idsChan <- id
		}

		close(idsChan)
	}()

	// here we should open cpu number of goroutines, but we know, the number of ids will not exceed 10000, so we are good
	for i := 0; i < len(ids); i++ {
		go c.BuildIssueItemConcurrent(issuesChan, idsChan, exitChan, add_attributes, re, auto_finalize_defect_type, auto_finalization_thredshold)
	}

	for i := 0; i < len(ids); i++ {
		<-exitChan
	}
	close(issuesChan)

	for issue := range issuesChan {
		issues = append(issues, issue)
	}

	return issues
}

// BuildIssueItemHelper method is a helper method for building
// the issue item struct.
func (c *RPConnector) BuildIssueItemHelper(id string, add_attributes bool, re bool, auto_finalize_defect_type bool, auto_finalization_threshold float32) IssueItem {
	logs := c.GetTestLog(id)
	// Make logs to string(in []byte format)
	log_after_marshal, _ := json.Marshal(logs)
	// This can be the input of GetPrediction
	testlog := string(log_after_marshal)

	var tfa_input common.TFAInput = c.BuildTFAInput(id, testlog, auto_finalize_defect_type, auto_finalization_threshold)
	prediction_json := c.GetPrediction(id, tfa_input)
	prediction := gjson.Get(prediction_json, "result.prediction").String()
	finalized_by_tfa := gjson.Get(prediction_json, "result.finalize").Bool()
	var issue_info IssueInfo = c.GetIssueInfoForSingleTestID(id)
	accuracy_score := gjson.Get(prediction_json, "result.probability_score").String()
	var issue_item IssueItem = IssueItem{Issue: issue_info, TestItemID: id}
	if finalized_by_tfa {
		if common.TFA_DEFECT_TYPE_TO_SUB_TYPE[prediction] != nil {

			prediction_code := common.TFA_DEFECT_TYPE[prediction]["locator"]
			issue_info.IssueType = prediction_code
		} else {
			log.Print("The predictions were not extracted correctly, so no update will be made!")
		}

		// Update the comment with re result
		if re {
			// test_item_detailed_info, _ := c.GetDetailedIssueInfoForSingleTestID(id)
			// test_item_name := gjson.Get(string(test_item_detailed_info), "content.0.name").String()
			// result, _ := json.Marshal(c.GGetDetailedIssueInfoForSingleTestIDetREResult(test_item_name))
			issue_info.Comment += c.GetREResult(id)
		}
		if add_attributes {
			prediction_name := common.TFA_DEFECT_TYPE[prediction]["longName"]
			err := c.updateAttributesForPrediction(id, prediction_name, accuracy_score, true)
			common.HandleError(err, "nopanic")
		}

		issue_item = IssueItem{Issue: issue_info, TestItemID: id}

	} else {
		// Added a default defect type
		if common.TFA_DEFECT_TYPE_TO_SUB_TYPE[prediction] != nil {

			prediction_code := common.TFA_DEFECT_TYPE_TO_SUB_TYPE[prediction]["locator"]
			// fmt.Println(prediction_code)
			issue_info.IssueType = prediction_code
		} else {
			log.Print("The predictions were not extracted correctly, so no update will be made!")
		}

		// Update the comment with re result
		if re {
			// test_item_detailed_info, _ := c.GetDetailedIssueInfoForSingleTestID(id)
			// test_item_name := gjson.Get(string(test_item_detailed_info), "content.0.name").String()
			// result, _ := json.Marshal(c.GGetDetailedIssueInfoForSingleTestIDetREResult(test_item_name))
			issue_info.Comment += c.GetREResult(id)
		}

		issue_item = IssueItem{Issue: issue_info, TestItemID: id}

		if add_attributes {
			prediction_name := common.TFA_DEFECT_TYPE_TO_SUB_TYPE[prediction]["longName"]
			err := c.updateAttributesForPrediction(id, prediction_name, accuracy_score, false)
			common.HandleError(err, "nopanic")
		}
	}
	return issue_item
}

// RERequestBody is the struct of request body for Recommendation Engine.
type RERequestBody struct {
	ProjectName string `json:"project_name"`
	LogMsg      string `json:"log_message"`
	TeamName    string `json:"TEAM"`
}

type LogMsgRequestBody struct {
	TestItemID []int  `json:"itemIds"`
	LogLevel   string `json:"logLevel"`
}
type REResults []REResult

// GetREResult can extract the returned re result.
func (c *RPConnector) GetREResult(id string) string {
	url := c.REURL
	method := http.MethodPost
	logMsg := c.GetTestLog(id)
	logMsgCombined := strings.Join(logMsg, "\n")
	var b RERequestBody = RERequestBody{ProjectName: c.ProjectName, LogMsg: logMsgCombined, TeamName: c.TeamName}

	var bb map[string]RERequestBody = map[string]RERequestBody{"data": b}

	body, _ := json.Marshal(bb)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, "", bytes.NewBuffer(body), c.Client)
	returned_ress := gjson.Get(string(data), "result").Str
	final_text := processREReturnedText(returned_ress)
	common.HandleError(err, "nopanic")

	return final_text
}

// processREReturnedText is a helper function which
// process the RE returned Information.
func processREReturnedText(re_result string) string {
	// re_result_links := strings.Split(re_result, " ")[2:13]
	// var link_map map[string]string = make(map[string]string)
	// for i, val := range re_result_links {
	// 	if i%4 == 0 {
	// 		link_map[val] = re_result_links[i+2]
	// 	}
	// }
	// final_text := "LINK, SCORE\n"
	// index := 1
	// for key, val := range link_map {
	// 	// [link1](http://foo.bar), O.5
	// 	final_text += fmt.Sprintf("[link%d](%s): %s\n", index, key, val)
	// 	index += 1
	// }

	return re_result
}

// BuildIssueItemConcurrent method builds Issue Item Concurrently.
func (c *RPConnector) BuildIssueItemConcurrent(issuesChan chan<- IssueItem, idsChan <-chan string, exitChan chan<- bool,
	add_attributes bool, re bool, auto_finalize_defect_type bool, auto_finalization_thredshold float32) {
	for {
		id, ok := <-idsChan
		if !ok {
			break
		}

		log.Printf("Getting prediction of test item(id): %s\n", id)
		issuesChan <- c.BuildIssueItemHelper(id, add_attributes, re, auto_finalize_defect_type, auto_finalization_thredshold)

	}
	exitChan <- true
}

// GetIssueInfoForSingleTestId method returns the issueinfo with the issue(test item) id.
func (c *RPConnector) GetDetailedIssueInfoForSingleTestID(id string) ([]byte, error) {
	if c.LaunchID == "" {
		c.LaunchID = c.GetLaunchID()
	}

	url := fmt.Sprintf("%s/api/v1/%s/item?filter.eq.id=%s&filter.eq.launchId=%s&isLatest=false&launchesLimit=0",
		c.RPURL, c.ProjectName, id, c.LaunchID)
	method := http.MethodGet
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)

	if err != nil {
		err = errors.Errorf("GetDetailedIssueInfoForSingleTestID error: %s", err)
	}

	return data, err
}

// GetIssueInfoForSingleTestId method returns the issueinfo with the issue(test item) id.
func (c *RPConnector) GetIssueInfoForSingleTestID(id string) IssueInfo {
	data, err := c.GetDetailedIssueInfoForSingleTestID(id)
	common.HandleError(err, "nopanic")

	issue_info_str := gjson.Get(string(data), "content.0.issue").String()

	var issue_info IssueInfo
	err = json.Unmarshal([]byte(issue_info_str), &issue_info)
	common.HandleError(err, "nopanic")

	return issue_info
}

// GetPrediction method returns the prediction extracted from the TFA Classifier.
func (c *RPConnector) GetPrediction(id string, tfa_input common.TFAInput) string {
	tfa_model := common.TFAModel{"data": tfa_input}

	model, err := json.Marshal(tfa_model)
	if err != nil {
		fmt.Println(err)
	}

	url := c.TFAURL
	method := http.MethodPost
	auth_token := c.AuthToken
	body := bytes.NewBuffer(model)

	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	if err != nil {
		common.HandleError(err, "nopanic")
	}

	return string(data)
}

// BuildTFAInput method builds the TFAInput struct with the test id and messages.
func (c *RPConnector) BuildTFAInput(test_id, messages string, auto_finalize_defect_type bool, auto_finalization_threshold float32) common.TFAInput {
	return common.TFAInput{ID: test_id, Project: c.ProjectName, Messages: messages, AutoFinalizeDefectType: auto_finalize_defect_type, FinalizationThreshold: auto_finalization_threshold}
}

func (c *RPConnector) GetAllTestInfos() map[string]string {
	if c.LaunchID == "" {
		c.LaunchID = c.GetLaunchID()
	}

	url := fmt.Sprintf("%s/api/v1/%s/item?filter.eq.issueType=ti001&filter.eq.launchId=%s&filter.eq."+
		"status=FAILED&isLatest=false&launchesLimit=0",
		c.RPURL, c.ProjectName, c.LaunchID)
	method := http.MethodGet
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	common.HandleError(err, "nopanic")

	a := gjson.Get(string(data), "content")

	var ret map[string]string = make(map[string]string)
	a.ForEach(func(_, m gjson.Result) bool {
		ret_id := m.Get("id").String()
		ret_name := m.Get("name").String()
		ret[ret_id] = ret_name
		return true
	})

	return ret
}

// GetAllTestIds returns all test ids from inside a test launch.
func (c *RPConnector) GetAllTestIds() []string {

	var ids []string
	for id, _ := range c.GetAllTestInfos() {
		ids = append(ids, id)
	}

	return ids
}

// GetTestLog returns the test log(test msg) for a test item.
func (c *RPConnector) GetTestLog(test_id string) []string {
	if c.LaunchID == "" {
		c.LaunchID = c.GetLaunchID()
	}

	url := fmt.Sprintf("%s/api/v1/%s/log?filter.eq.item=%s&filter.eq.launchId=%s",
		c.RPURL, c.ProjectName, test_id, c.LaunchID)
	method := http.MethodGet
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	common.HandleError(err, "nopanic")

	a := gjson.Get(string(data), "content")

	var ret []string

	a.ForEach(func(_, m gjson.Result) bool {
		ret = append(ret, m.Get("message").String())

		return true
	})

	return ret
}

type attribute map[string]string

func (c *RPConnector) getExistingAtrributeByID(id string) Attributes {
	url := fmt.Sprintf("%s/api/v1/%s/item/uuid/%s", c.RPURL, c.ProjectName, id)
	method := http.MethodGet
	auth_token := c.AuthToken

	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	common.HandleError(err, "nopanic")

	if err != nil {
		err = fmt.Errorf("Get attibute failed:%w", err)
		common.HandleError(err, "nopanic")
	}

	attrs := gjson.Get(string(data), "attributes").String()
	attr := []attribute{}

	err = json.Unmarshal([]byte(attrs), &attr)
	common.HandleError(err, "nopanic")

	return Attributes{"attributes": attr}
}

func (c *RPConnector) updateAttributesForPrediction(id, prediction, accuracy_score string, finalized_by_tfa bool) error {
	existingAttribute := c.getExistingAtrributeByID(id)
	tfa_prediction_attr := attribute{

		"key":   "AI Prediction",
		"value": prediction,
	}
	finalized_by_tfa_attr := attribute{

		"key":   "Finalized_By",
		"value": "TFA",
	}
	if finalized_by_tfa {
		existingAttribute["attributes"] = append(existingAttribute["attibutes"], finalized_by_tfa_attr)
	}
	var tfa_accuracy_score attribute
	if accuracy_score != "" {
		tfa_accuracy_score = attribute{

			"key":   "Prediction Score",
			"value": accuracy_score,
		}
	}
	existingAttribute["attributes"] = append(existingAttribute["attributes"], tfa_prediction_attr)
	if accuracy_score != "" {
		existingAttribute["attributes"] = append(existingAttribute["attributes"], tfa_accuracy_score)
	}

	url := fmt.Sprintf("%s/api/v1/%s/item/%s/update", c.RPURL, c.ProjectName, id)
	method := http.MethodPut
	auth_token := c.AuthToken
	d, err := json.Marshal(existingAttribute)
	common.HandleError(err, "nopanic")

	body := bytes.NewBuffer(d)
	_, _, err = common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)

	if err != nil {
		err = fmt.Errorf("updated attibute failed:%w", err)
	}

	log.Printf("Updated the test item(id): %s with it's prediction %s\n", id, prediction)

	return err
}

func getExistingDefectTypeLocatorID(gjson_obj []gjson.Result, defect_type string) (string, bool) {
	for _, v := range gjson_obj {
		defect_type_info := v.Map()
		if defect_type_info["longName"].String() == defect_type {
			return defect_type_info["locator"].String(), true
		}
	}

	return "", false
}

// InitConnector create defect types before doing all sync/update job
// this method will run before everything.
func (c *RPConnector) InitConnector() {
	c.RPURL = strings.TrimSuffix(c.RPURL, "/")
	fmt.Println("Initializing Defect Types...")
	url := fmt.Sprintf("%s/api/v1/%s/settings", c.RPURL, c.ProjectName)
	method := http.MethodGet
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, success, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	common.HandleError(err, "panic")

	if !success {
		fmt.Printf("created defect types failed, please use superadmin auth_token%t", success)

		return
	}

	ti_sub := gjson.Get(string(data), "subTypes.TO_INVESTIGATE").Array()

	for _, sub_type := range common.PREDICTED_SUB_TYPES {
		locator, ok := getExistingDefectTypeLocatorID(ti_sub, sub_type["longName"])
		if !ok {
			d, _ := json.Marshal(sub_type)
			url := fmt.Sprintf("%s/api/v1/%s/settings/sub-type", c.RPURL, c.ProjectName)
			method := http.MethodPost
			auth_token := c.AuthToken
			body := bytes.NewBuffer(d)

			data, success, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
			if err != nil {

				err = errors.Errorf("read response body failed: %s", err)
				common.HandleError(err, "nopanic")
			}

			if !success {
				fmt.Printf("created defect types failed: %t\n", success)

				return
			}

			sub_type["locator"] = gjson.Get(string(data), "locator").String()
		} else {
			sub_type["locator"] = locator
		}
	}
}

func (c *RPConnector) GetLaunchIDByName() string {
	launchinfo := strings.Split(c.LaunchName, "#")
	launchName := url.QueryEscape(strings.TrimSpace(launchinfo[0]))
	url := fmt.Sprintf("%s/api/v1/%s/launch?filter.eq.name=%s&filter.eq.number=%s",
		c.RPURL, c.ProjectName, launchName, launchinfo[1])
	method := http.MethodGet
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	if err != nil {
		fmt.Printf("Get Launch Id failed: %v", string(data))
	}

	launchId := gjson.Get(string(data), "content.0.id").String()

	return launchId

}

func (c *RPConnector) GetLaunchIDByUUID() string {

	url := fmt.Sprintf("%s/api/v1/%s/launch?filter.eq.uuid=%s",
		c.RPURL, c.ProjectName, c.LaunchUUID)
	method := http.MethodGet
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	if err != nil {
		fmt.Printf("Get Launch Id failed: %v", string(data))
	}

	launchId := gjson.Get(string(data), "content.0.id").String()

	return launchId

}

// GetLaunchID returns launch id with the launch name
// this method will be called when user don't have launchid input.
func (c *RPConnector) GetLaunchID() string {
	if c.LaunchName != "" {
		return c.GetLaunchIDByName()
	} else if c.LaunchUUID != "" {
		return c.GetLaunchIDByUUID()
	} else {
		return ""
	}
}

type Attributes map[string][]attribute

func (c *RPConnector) GetAttributesByID(id string) Attributes {
	url := fmt.Sprintf("%s/api/v1/%s/item/uuid/%s",
		c.RPURL, c.ProjectName, id)
	method := http.MethodGet
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, auth_token, body, c.Client)
	if err != nil {
		fmt.Printf("Get Launch Id failed: %v", string(data))
	}

	attributes := gjson.Get(string(data), "attributes").String()

	attr := Attributes{}
	err = json.Unmarshal([]byte(attributes), &attr)
	common.HandleError(err, "nopanic")

	return attr
}

// Revert function will set all test items back to the
// Original defect type.
func (c *RPConnector) RevertUpdatedList(verbose bool) common.GeneralUpdatedList {
	ids := c.GetAllTestIds()
	issues := Issues{}

	for _, id := range ids {
		issues = append(issues, c.revertHelper(id))
	}

	return UpdatedList{IssuesList: issues}
}

func (c *RPConnector) revertHelper(id string) IssueItem {
	var issue_info IssueInfo = c.GetIssueInfoForSingleTestID(id)
	issue_info.IssueType = "ti001"

	var issue_item IssueItem = IssueItem{Issue: issue_info, TestItemID: id}

	return issue_item
}
