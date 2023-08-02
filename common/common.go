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

var TFA_DEFECT_TYPE_TO_SUB_TYPE map[string]PREDICTED_SUB_TYPE = map[string]PREDICTED_SUB_TYPE{
	"Automation Bug": PREDICTED_AUTOMATION_BUG,
	"Product Bug":    PREDICTED_PRODUCT_BUG,
	"System Issue":   PREDICTED_SYSTEM_BUG,
	"No Defect":      PREDICTED_NO_DEFECT,
}
var TFA_DEFECT_TYPE map[string]MAIN_DEFECT_TYPE = map[string]MAIN_DEFECT_TYPE{
	"Automation Bug": AUTOMATION_BUG,
	"Product Bug":    PRODUCT_BUG,
	"System Issue":   SYSTEM_BUG,
	"No Defect":      NO_DEFECT,
}

type PREDICTED_SUB_TYPE map[string]string
type MAIN_DEFECT_TYPE map[string]string

var PREDICTED_SUB_TYPES map[string]PREDICTED_SUB_TYPE = map[string]PREDICTED_SUB_TYPE{
	"PREDICTED_AUTOMATION_BUG": PREDICTED_AUTOMATION_BUG,
	"PREDICTED_SYSTEM_BUG":     PREDICTED_SYSTEM_BUG,
	"PREDICTED_PRODUCT_BUG":    PREDICTED_PRODUCT_BUG,
	"PREDICTED_NO_DEFECT":      PREDICTED_NO_DEFECT,
}
var MAIN_DEFECT_TYPES map[string]MAIN_DEFECT_TYPE = map[string]MAIN_DEFECT_TYPE{
	"AUTOMATION_BUG": AUTOMATION_BUG,
	"SYSTEM_BUG":     SYSTEM_BUG,
	"PRODUCT_BUG":    PRODUCT_BUG,
	"NO_DEFECT":      NO_DEFECT,
}

var NO_DEFECT MAIN_DEFECT_TYPE = MAIN_DEFECT_TYPE{
	"locator":   "nd001",
	"typeRef":   "NO_DEFECT",
	"longName":  "No Defect",
	"shortName": "ND",
	"color":     "#777777",
}
var AUTOMATION_BUG MAIN_DEFECT_TYPE = MAIN_DEFECT_TYPE{
	"locator":   "ab001",
	"typeRef":   "AUTOMATION_BUG",
	"longName":  "Automation Bug",
	"shortName": "AB",
	"color":     "#f7d63e",
}
var PRODUCT_BUG MAIN_DEFECT_TYPE = MAIN_DEFECT_TYPE{
	"locator":   "pb001",
	"typeRef":   "PRODUCT_BUG",
	"longName":  "Product Bug",
	"shortName": "PB",
	"color":     "#ec3900",
}

var SYSTEM_BUG MAIN_DEFECT_TYPE = MAIN_DEFECT_TYPE{
	"locator":   "si001",
	"typeRef":   "SYSTEM_ISSUE",
	"longName":  "System Issue",
	"shortName": "SI",
	"color":     "#0274d1",
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
