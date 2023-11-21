// Package connectors is the package for all connectors struct and its
// coressponding methods
package connectors

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/RedHatQE/tfacon/common"
	"github.com/tidwall/gjson"
)

func TestRPConnector_updateAttributesForPrediction(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		id               string
		prediction       string
		accurracy_score  string
		finalized_by_tfa bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			if err := c.updateAttributesForPrediction(tt.args.id, tt.args.prediction, tt.args.accurracy_score, tt.args.finalized_by_tfa); (err != nil) != tt.wantErr {
				t.Errorf("RPConnector.updateAttributesForPrediction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getExistingDefectTypeLocatorID(t *testing.T) {
	type args struct {
		gjson_obj   []gjson.Result
		defect_type string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getExistingDefectTypeLocatorID(tt.args.gjson_obj, tt.args.defect_type)
			if got != tt.want {
				t.Errorf("getExistingDefectTypeLocatorID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getExistingDefectTypeLocatorID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRPConnector_InitConnector(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			c.InitConnector()
		})
	}
}

func TestRPConnector_GetLaunchID(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			if got := c.GetLaunchID(); got != tt.want {
				t.Errorf("RPConnector.GetLaunchID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_Validate(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		verbose bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			got, err := c.Validate(tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("RPConnector.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RPConnector.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_validateTFAURL(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		verbose bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			got, err := c.validateTFAURL(tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("RPConnector.validateTFAURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RPConnector.validateTFAURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_validateLaunchInfo(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		verbose bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			got, err := c.validateLaunchInfo(tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("RPConnector.validateLaunchInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RPConnector.validateLaunchInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_validateProjectName(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		verbose bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			got, err := c.validateProjectName(tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("RPConnector.validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RPConnector.validateProjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_validateRPURLAndAuthToken(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		verbose bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			got, err := c.validateRPURLAndAuthToken(tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("RPConnector.validateRPURLAndAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RPConnector.validateRPURLAndAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_String(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("RPConnector.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_UpdateAll(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		updatedListOfIssues common.GeneralUpdatedList
		verbose             bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			c.UpdateAll(tt.args.updatedListOfIssues, tt.args.verbose)
		})
	}
}

// func TestRPConnector_BuildUpdatedList(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	type args struct {
// 		ids            []string
// 		concurrent     bool
// 		add_attributes bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   common.GeneralUpdatedList
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			if got := c.BuildUpdatedList(tt.args.ids, tt.args.concurrent, tt.args.add_attributes); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RPConnector.BuildUpdatedList() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestRPConnector_BuildIssues(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	type args struct {
// 		ids            []string
// 		concurrent     bool
// 		add_attributes bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   Issues
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			if got := c.BuildIssues(tt.args.ids, tt.args.concurrent, tt.args.add_attributes); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RPConnector.BuildIssues() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestRPConnector_BuildIssuesConcurrent(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	type args struct {
// 		ids            []string
// 		add_attributes bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   Issues
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			if got := c.BuildIssuesConcurrent(tt.args.ids, tt.args.add_attributes); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RPConnector.BuildIssuesConcurrent() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestRPConnector_BuildIssueItemHelper(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	type args struct {
// 		id             string
// 		add_attributes bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   IssueItem
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			if got := c.BuildIssueItemHelper(tt.args.id, tt.args.add_attributes); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RPConnector.BuildIssueItemHelper() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestRPConnector_BuildIssueItemConcurrent(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	type args struct {
// 		issuesChan     chan<- IssueItem
// 		idsChan        <-chan string
// 		exitChan       chan<- bool
// 		add_attributes bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			c.BuildIssueItemConcurrent(tt.args.issuesChan, tt.args.idsChan, tt.args.exitChan, tt.args.add_attributes)
// 		})
// 	}
// }

// func TestRPConnector_GetIssueInfoForSingleTestID(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	type args struct {
// 		id string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   IssueInfo
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			if got := c.GetIssueInfoForSingleTestID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RPConnector.GetIssueInfoForSingleTestID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestRPConnector_GetPrediction(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		id        string
		tfa_input common.TFAInput
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			if got := c.GetPrediction(tt.args.id, tt.args.tfa_input); got != tt.want {
				t.Errorf("RPConnector.GetPrediction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPConnector_BuildTFAInput(t *testing.T) {
	type fields struct {
		LaunchID    string
		LaunchName  string
		ProjectName string
		AuthToken   string
		RPURL       string
		Client      *http.Client
		TFAURL      string
	}
	type args struct {
		test_id                     string
		messages                    string
		auto_finalize_defect_type   bool
		auto_finalization_threshold float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.TFAInput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPConnector{
				LaunchID:    tt.fields.LaunchID,
				LaunchName:  tt.fields.LaunchName,
				ProjectName: tt.fields.ProjectName,
				AuthToken:   tt.fields.AuthToken,
				RPURL:       tt.fields.RPURL,
				Client:      tt.fields.Client,
				TFAURL:      tt.fields.TFAURL,
			}
			if got := c.BuildTFAInput(tt.args.test_id, tt.args.messages, tt.args.auto_finalize_defect_type, tt.args.auto_finalization_threshold); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RPConnector.BuildTFAInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestRPConnector_GetAllTestIds(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   []string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			if got := c.GetAllTestIds(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RPConnector.GetAllTestIds() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestRPConnector_GetTestLog(t *testing.T) {
// 	type fields struct {
// 		LaunchID    string
// 		LaunchName  string
// 		ProjectName string
// 		AuthToken   string
// 		RPURL       string
// 		Client      *http.Client
// 		TFAURL      string
// 	}
// 	type args struct {
// 		test_id string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   []string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RPConnector{
// 				LaunchID:    tt.fields.LaunchID,
// 				LaunchName:  tt.fields.LaunchName,
// 				ProjectName: tt.fields.ProjectName,
// 				AuthToken:   tt.fields.AuthToken,
// 				RPURL:       tt.fields.RPURL,
// 				Client:      tt.fields.Client,
// 				TFAURL:      tt.fields.TFAURL,
// 			}
// 			if got := c.GetTestLog(tt.args.test_id); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RPConnector.GetTestLog() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
