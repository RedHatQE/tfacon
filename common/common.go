// Package common contains all shared structs(data structures) required for all connectors
package common

// GeneralUpdatedList is an updated list of object, with the prediction from TFA classifier
// each connector should have it's own UpdatedList structure and implement the
// GeneralUpdatedList interface.
type GeneralUpdatedList interface {
	GetSelf() GeneralUpdatedList
}

type (
	// TFAModel is the data structure for describing
	// the request body for TFA Classifer.
	TFAModel map[string]TFAInput
	// TFAInput is the data structure for describing
	// the three input params for TFA Classifier.
	TFAInput struct {
		ID       string `json:"id"`
		Project  string `json:"project"`
		Messages string `json:"messages"`
	}
)

const DEFAULT_LOG_PATH string = "./tfacon.log"

var TFA_DEFECT_TYPE_TO_SUB_TYPE map[string]PREDICTED_SUB_TYPE = map[string]PREDICTED_SUB_TYPE{
	"Automation Bug": PREDICTED_AUTOMATION_BUG,
	"Product Bug":    PREDICTED_PRODUCT_BUG,
	"System Issue":   PREDICTED_SYSTEM_BUG,
	"No Defect":      PREDICTED_NO_DEFECT,
}

type PREDICTED_SUB_TYPE map[string]string

var PREDICTED_SUB_TYPES map[string]PREDICTED_SUB_TYPE = map[string]PREDICTED_SUB_TYPE{
	"PREDICTED_AUTOMATION_BUG": PREDICTED_AUTOMATION_BUG,
	"PREDICTED_SYSTEM_BUG":     PREDICTED_SYSTEM_BUG,
	"PREDICTED_PRODUCT_BUG":    PREDICTED_PRODUCT_BUG,
	"PREDICTED_NO_DEFECT":      PREDICTED_NO_DEFECT,
}

var PREDICTED_AUTOMATION_BUG PREDICTED_SUB_TYPE = PREDICTED_SUB_TYPE{
	"typeRef":   "TO_INVESTIGATE",
	"longName":  "Predicted Automation Bug",
	"shortName": "TIA",
	"color":     "#ffeeaa",
}

var PREDICTED_SYSTEM_BUG PREDICTED_SUB_TYPE = PREDICTED_SUB_TYPE{
	"typeRef":   "TO_INVESTIGATE",
	"longName":  "Predicted System Issue",
	"shortName": "TIS",
	"color":     "#aaaaff",
}

var PREDICTED_PRODUCT_BUG PREDICTED_SUB_TYPE = PREDICTED_SUB_TYPE{
	"typeRef":   "TO_INVESTIGATE",
	"longName":  "Predicted Product Bug",
	"shortName": "TIP",
	"color":     "#ffaaaa",
}

var PREDICTED_NO_DEFECT PREDICTED_SUB_TYPE = PREDICTED_SUB_TYPE{
	"typeRef":   "TO_INVESTIGATE",
	"longName":  "Predicted No Defect",
	"shortName": "TND",
	"color":     "#C1BCB4",
}
