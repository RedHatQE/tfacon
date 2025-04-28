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
	"strconv"
	"strings"

	"github.com/RedHatQE/tfacon/common"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Issues is a list of IssueItem.
type Issues []IssueItem

// REResult is the RE request return result struct
type REResult struct {
	LaunchID     string `json:"launch_id"`
	LaunchName   string `json:"launch_name"`
	TestItemID   string `json:"test_item_id"`
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
func (c *RPConnector) String() string {
	str := common.StructToString(c)
	return str
}

// UpdateAll method is the interface method for tfacon interface
// it update the test items with the built updated_list.
func (c *RPConnector) UpdateAll(updatedListOfIssues common.GeneralUpdatedList, verbose bool) {
	if len(updatedListOfIssues.GetSelf().(UpdatedList).IssuesList) == 0 {
		return
	}

	jsonUpdatedListOfIssues, _ := json.Marshal(updatedListOfIssues)

	log.Println("Updating All Test Items With Predictions")

	url := fmt.Sprintf("%s/api/v1/%s/item", c.RPURL, c.ProjectName)
	method := "PUT"
	authToken := c.AuthToken
	body := bytes.NewBuffer(jsonUpdatedListOfIssues)

	data, success, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
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
func (c *RPConnector) BuildUpdatedList(ids []string, tfaConfig common.TFAConfig) common.GeneralUpdatedList {
	return UpdatedList{IssuesList: c.BuildIssues(ids, tfaConfig)}
}

// BuildIssues method build the issue struct.
func (c *RPConnector) BuildIssues(ids []string, tfaConfig common.TFAConfig) Issues {
	issues := Issues{}

	if tfaConfig.Concurrent {
		return c.BuildIssuesConcurrent(ids, tfaConfig)
	}

	for _, id := range ids {
		log.Printf("Getting prediction of test item(id): %s\n", id)
		issues = append(issues, c.BuildIssueItemHelper(id, tfaConfig))
	}

	return issues
}

// BuildIssuesConcurrent methods builds the issues struct concurrently.
func (c *RPConnector) BuildIssuesConcurrent(ids []string, tfaConfig common.TFAConfig) Issues {
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
		go c.BuildIssueItemConcurrent(issuesChan, idsChan, exitChan, tfaConfig)
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
func (c *RPConnector) BuildIssueItemHelper(id string, tfaConfig common.TFAConfig) IssueItem {
	logs := c.GetTestLog(id)
	// Make logs to string(in []byte format)
	logAfterMarshal, _ := json.Marshal(logs)
	// This can be the input of GetPrediction
	testlog := string(logAfterMarshal)

	var tfaInput common.TFAInput = c.BuildTFAInput(id, testlog, tfaConfig.AutoFinalizeDefectType, tfaConfig.AutoFinalizationThreshold)
	predictionJSON := c.GetPrediction(id, tfaInput)
	prediction := gjson.Get(predictionJSON, "result.prediction").String()
	finalizedByTfa := gjson.Get(predictionJSON, "result.finalize").Bool()
	var issueInfo = c.GetIssueInfoForSingleTestID(id)
	accuracyScore := gjson.Get(predictionJSON, "result.probability_score").String()
	var issueItem = IssueItem{Issue: issueInfo, TestItemID: id}
	if finalizedByTfa {
		if common.TFADefectTypeToSubType[prediction] != nil {
			predictionCode := common.TFADefectType[prediction]["locator"]
			issueInfo.IssueType = predictionCode
			log.Printf("Finalizer prediction code is: %s\n", predictionCode)
			issueInfo.Comment = "### TFA-C Auto Finalized\n" + issueInfo.Comment
		} else {
			log.Println("The predictions were not extracted correctly, so no finalizer update will be made!")
		}

		// Update the comment with re result
		if tfaConfig.Re {
			if issueInfo.Comment != "" {
				issueInfo.Comment += "\n"
				issueInfo.Comment += c.GetREResult(id)
			} else {
				issueInfo.Comment += c.GetREResult(id)
			}
		}
		if tfaConfig.AddAttributes {
			predictionName := common.TFADefectType[prediction]["longName"]
			err := c.updateAttributesForPrediction(id, predictionName, accuracyScore, true)
			common.HandleError(err, "nopanic")
		}

		issueItem = IssueItem{Issue: issueInfo, TestItemID: id}

	} else {
		// Added a default defect type
		if common.TFADefectTypeToSubType[prediction] != nil {

			predictionCode := common.TFADefectTypeToSubType[prediction]["locator"]
			// fmt.Println(predictionCode)
			issueInfo.IssueType = predictionCode
		} else {
			log.Println("The predictions were not extracted correctly, so no update will be made!")
		}

		// Update the comment with re result
		if tfaConfig.Re {
			if issueInfo.Comment != "" {
				issueInfo.Comment += "\n"
				issueInfo.Comment += c.GetREResult(id)
			} else {
				issueInfo.Comment += c.GetREResult(id)
			}
		}

		issueItem = IssueItem{Issue: issueInfo, TestItemID: id}

		if tfaConfig.AddAttributes {
			predictionName := common.TFADefectTypeToSubType[prediction]["longName"]
			err := c.updateAttributesForPrediction(id, predictionName, accuracyScore, false)
			common.HandleError(err, "nopanic")
		}
	}
	return issueItem
}

// RERequestBody is the struct of request body for Recommendation Engine.
type RERequestBody struct {
	ProjectName string `json:"project_name"`
	LogMsg      string `json:"log_message"`
	TeamName    string `json:"TEAM"`
}

// LogMsgRequestBody is the struct for log message body
type LogMsgRequestBody struct {
	TestItemID []int  `json:"itemIds"`
	LogLevel   string `json:"logLevel"`
}

// REResults is list of REResult
type REResults []REResult

// GetREResult can extract the returned re result.
func (c *RPConnector) GetREResult(id string) string {
	url := c.REURL
	method := http.MethodPost
	logMsg := c.GetTestLog(id)
	logMsgCombined := strings.Join(logMsg, "\n")
	var b = RERequestBody{ProjectName: c.ProjectName, LogMsg: logMsgCombined, TeamName: c.TeamName}

	var bb = map[string]RERequestBody{"data": b}

	body, _ := json.Marshal(bb)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, "", bytes.NewBuffer(body), c.Client)
	returnedRess := gjson.Get(string(data), "result").Str
	common.HandleError(err, "nopanic")
	finalText, err := processREReturnedText(returnedRess)
	common.HandleError(err, "nopanic")

	return finalText
}

// processREReturnedText is a helper function which
// process the RE returned Information.
func processREReturnedText(reResult string) (string, error) {
	var errRet error
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("reResult hasn't been retrieved correctly")
			errRet = errors.Errorf("reResult hasn't been retrieved correctly: %s", r)
			return
		}
	}()
	if reResult != "No Similar Test Case Found" {
		// if reResult is No Similar Test Case Found, we don't need to do preprocessing
		var allStringsFromREResult []string = strings.Split(reResult, ",")
		// finalStrings stores all valid non-empty strings from the reResult one by one in array
		finalStrings := make([]string, 0)
		for _, str := range allStringsFromREResult {
			currentSplitRes := strings.Split(str, " ")
			for _, str := range currentSplitRes {
				if str != "" {
					finalStrings = append(finalStrings, str)
				}
			}
		}
		finalStrings = finalStrings[1:]
		var finalInfo = make(map[string][]string)
		resultNumber := 1
		for i := 0; i < len(finalStrings); i += 2 {
			keyName := "Result " + strconv.Itoa(resultNumber)
			finalInfo[keyName] = []string{finalStrings[i], finalStrings[i+1]}
			resultNumber++

		}
		finalText := "| TFA-R | Similarity Score |\n|--|--|\n"
		for key, val := range finalInfo {
			// [link1](http://foo.bar), O.5
			finalText += fmt.Sprintf("|[%s](%s) | %s |\n", key, val[0], val[1])
		}

		reResult = finalText
	}

	return reResult, errRet
}

// BuildIssueItemConcurrent method builds Issue Item Concurrently.
func (c *RPConnector) BuildIssueItemConcurrent(issuesChan chan<- IssueItem, idsChan <-chan string, exitChan chan<- bool,
	tfaConfig common.TFAConfig) {
	for {
		id, ok := <-idsChan
		if !ok {
			break
		}

		log.Printf("Getting prediction of test item(id): %s\n", id)
		issuesChan <- c.BuildIssueItemHelper(id, tfaConfig)

	}
	exitChan <- true
}

// GetDetailedIssueInfoForSingleTestID method returns the issueinfo with the issue(test item) id.
func (c *RPConnector) GetDetailedIssueInfoForSingleTestID(id string) ([]byte, error) {
	if c.LaunchID == "" {
		c.LaunchID = c.GetLaunchID()
	}

	url := fmt.Sprintf("%s/api/v1/%s/item?filter.eq.id=%s&filter.eq.launchId=%s&isLatest=false&launchesLimit=0",
		c.RPURL, c.ProjectName, id, c.LaunchID)
	method := http.MethodGet
	authToken := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)

	if err != nil {
		err = errors.Errorf("GetDetailedIssueInfoForSingleTestID error: %s", err)
	}

	return data, err
}

// GetIssueInfoForSingleTestID method returns the issueinfo with the issue(test item) id.
func (c *RPConnector) GetIssueInfoForSingleTestID(id string) IssueInfo {
	data, err := c.GetDetailedIssueInfoForSingleTestID(id)
	common.HandleError(err, "nopanic")

	issueInfoStr := gjson.Get(string(data), "content.0.issue").String()

	var issueInfo IssueInfo
	err = json.Unmarshal([]byte(issueInfoStr), &issueInfo)
	common.HandleError(err, "nopanic")

	return issueInfo
}

// GetPrediction method returns the prediction extracted from the TFA Classifier.
func (c *RPConnector) GetPrediction(id string, tfaInput common.TFAInput) string {
	tfaModel := common.TFAModel{"data": tfaInput}

	model, err := json.Marshal(tfaModel)
	if err != nil {
		fmt.Println(err)
	}

	url := c.TFAURL
	method := http.MethodPost
	authToken := c.AuthToken
	body := bytes.NewBuffer(model)

	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
	if err != nil {
		common.HandleError(err, "nopanic")
	}

	return string(data)
}

// BuildTFAInput method builds the TFAInput struct with the test id and messages.
func (c *RPConnector) BuildTFAInput(testID, messages string, autoFinalizeDefectType bool, autoFinalizationThreshold float32) common.TFAInput {
	return common.TFAInput{ID: testID, Project: c.ProjectName, Messages: messages, AutoFinalizeDefectType: autoFinalizeDefectType, FinalizationThreshold: autoFinalizationThreshold}
}

// GetAllTestInfos send rest query to RP and return failed result of given issueType
// return all failed items is issueType is empty
func (c *RPConnector) GetAllTestInfos(issueType string, pageSize int) map[string]string {
	// default filter type 'ti001' as 'to be investigated'
	var url string
	if c.LaunchID == "" {
		c.LaunchID = c.GetLaunchID()
	}

	if issueType == "" {
		url = fmt.Sprintf("%s/api/v1/%s/item?filter.eq.hasChildren=false&filter.eq.launchId=%s&filter.eq."+
			"status=FAILED&isLatest=false&launchesLimit=0&page.size=%d",
			c.RPURL, c.ProjectName, c.LaunchID, pageSize)
	} else {
		url = fmt.Sprintf("%s/api/v1/%s/item?filter.eq.hasChildren=false&filter.eq.issueType=%s&filter.eq.launchId=%s&filter.eq."+
			"status=FAILED&isLatest=false&launchesLimit=0&page.size=%d",
			c.RPURL, c.ProjectName, issueType, c.LaunchID, pageSize)
	}

	method := http.MethodGet
	authToken := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
	common.HandleError(err, "nopanic")

	a := gjson.Get(string(data), "content")

	var ret = make(map[string]string)
	a.ForEach(func(_, m gjson.Result) bool {
		retID := m.Get("id").String()
		retName := m.Get("name").String()
		ret[retID] = retName
		return true
	})

	return ret
}

// GetAllTestIds returns all test ids from inside a test launch.
func (c *RPConnector) GetAllTestIds() []string {
	var ids []string
	issueType := "ti001"
	pageSize := 20
	for id := range c.GetAllTestInfos(issueType, pageSize) {
		ids = append(ids, id)
	}

	return ids
}

// GetTestLog returns the test log(test msg) for a test item.
func (c *RPConnector) GetTestLog(testID string) []string {
	if c.LaunchID == "" {
		c.LaunchID = c.GetLaunchID()
	}

	url := fmt.Sprintf("%s/api/v1/%s/log?filter.eq.item=%s&filter.eq.launchId=%s",
		c.RPURL, c.ProjectName, testID, c.LaunchID)
	method := http.MethodGet
	authToken := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
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
	authToken := c.AuthToken

	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
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

func (c *RPConnector) updateAttributesForPrediction(id, prediction, accuracyScore string, finalizedByTfa bool) error {
	existingAttribute := c.getExistingAtrributeByID(id)
	tfaPredictionAttr := attribute{

		"key":   "AI Prediction",
		"value": prediction,
	}
	finalizedByTfaAttr := attribute{

		"key":   "Finalized_By",
		"value": "TFA",
	}
	if finalizedByTfa {
		existingAttribute["attributes"] = append(existingAttribute["attibutes"], finalizedByTfaAttr)
	}
	var tfaAccuracyScore attribute
	if accuracyScore != "" {
		tfaAccuracyScore = attribute{

			"key":   "Prediction Score",
			"value": accuracyScore,
		}
	}
	existingAttribute["attributes"] = append(existingAttribute["attributes"], tfaPredictionAttr)
	if accuracyScore != "" {
		existingAttribute["attributes"] = append(existingAttribute["attributes"], tfaAccuracyScore)
	}

	url := fmt.Sprintf("%s/api/v1/%s/item/%s/update", c.RPURL, c.ProjectName, id)
	method := http.MethodPut
	authToken := c.AuthToken
	d, err := json.Marshal(existingAttribute)
	common.HandleError(err, "nopanic")

	body := bytes.NewBuffer(d)
	_, _, err = common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)

	if err != nil {
		err = fmt.Errorf("updated attibute failed:%w", err)
	}

	log.Printf("Updated the test item(id): %s attributes with it's prediction %s\n", id, prediction)

	return err
}

func getExistingDefectTypeLocatorID(gjsonObj []gjson.Result, defectType string) (string, bool) {
	for _, v := range gjsonObj {
		defectTypeInfo := v.Map()
		if defectTypeInfo["longName"].String() == defectType {
			return defectTypeInfo["locator"].String(), true
		}
	}

	return "", false
}

// InitConnector create defect types before doing all sync/update job
// this method will run before everything.
func (c *RPConnector) InitConnector() {
	fmt.Println("Initializing Defect Types...")
	url := fmt.Sprintf("%s/api/v1/%s/settings", c.RPURL, c.ProjectName)
	method := http.MethodGet
	authToken := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, success, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
	common.HandleError(err, "panic")

	if !success {
		fmt.Printf("created defect types failed, please use superadmin authToken%t", success)

		return
	}

	tiSub := gjson.Get(string(data), "subTypes.TO_INVESTIGATE").Array()

	for _, subType := range common.PredictedSubTypes {
		locator, ok := getExistingDefectTypeLocatorID(tiSub, subType["longName"])
		if !ok {
			d, _ := json.Marshal(subType)
			url := fmt.Sprintf("%s/api/v1/%s/settings/sub-type", c.RPURL, c.ProjectName)
			method := http.MethodPost
			authToken := c.AuthToken
			body := bytes.NewBuffer(d)

			data, success, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
			if err != nil {

				err = errors.Errorf("read response body failed: %s", err)
				common.HandleError(err, "nopanic")
			}

			if !success {
				fmt.Printf("created defect types failed: %t\n", success)

				return
			}

			subType["locator"] = gjson.Get(string(data), "locator").String()
		} else {
			subType["locator"] = locator
		}
	}
}

// GetLaunchIDByName send rest query to RP with launch name and return the launch ID
func (c *RPConnector) GetLaunchIDByName() string {
	launchinfo := strings.Split(c.LaunchName, "#")
	launchName := url.QueryEscape(strings.TrimSpace(launchinfo[0]))
	url := fmt.Sprintf("%s/api/v1/%s/launch?filter.eq.name=%s&filter.eq.number=%s",
		c.RPURL, c.ProjectName, launchName, launchinfo[1])
	method := http.MethodGet
	authToken := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
	if err != nil {
		fmt.Printf("Get Launch Id failed: %v", string(data))
	}

	launchID := gjson.Get(string(data), "content.0.id").String()

	return launchID

}

// GetLaunchIDByUUID send rest query to RP with launch uuid and return the launch ID
func (c *RPConnector) GetLaunchIDByUUID() string {

	url := fmt.Sprintf("%s/api/v1/%s/launch?filter.eq.uuid=%s",
		c.RPURL, c.ProjectName, c.LaunchUUID)
	method := http.MethodGet
	authToken := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
	if err != nil {
		fmt.Printf("Get Launch Id failed: %v", string(data))
	}

	launchID := gjson.Get(string(data), "content.0.id").String()

	return launchID

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

// Attributes map of attribute
type Attributes map[string][]attribute

// GetAttributesByID send rest query to RP and return attributes
func (c *RPConnector) GetAttributesByID(id string) Attributes {
	url := fmt.Sprintf("%s/api/v1/%s/item/uuid/%s",
		c.RPURL, c.ProjectName, id)
	method := http.MethodGet
	authToken := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, _, err := common.SendHTTPRequest(context.Background(), method, url, authToken, body, c.Client)
	if err != nil {
		fmt.Printf("Get Launch Id failed: %v", string(data))
	}

	attributes := gjson.Get(string(data), "attributes").String()

	attr := Attributes{}
	err = json.Unmarshal([]byte(attributes), &attr)
	common.HandleError(err, "nopanic")

	return attr
}

// RevertUpdatedList function will set all test items back to the
// Original defect type.
func (c *RPConnector) RevertUpdatedList(verbose bool) common.GeneralUpdatedList {
	ids := c.GetAllTestIds()
	issues := Issues{}

	for _, id := range ids {
		issues = append(issues, c.revertHelper(id))
	}

	return UpdatedList{IssuesList: issues}
}

// revertHelper return the IssueItem with default issue type "ti001"
func (c *RPConnector) revertHelper(id string) IssueItem {
	var issueInfo = c.GetIssueInfoForSingleTestID(id)
	issueInfo.IssueType = "ti001"

	var issueItem = IssueItem{Issue: issueInfo, TestItemID: id}

	return issueItem
}
