package gaia

type CommonResponse struct {
	Result  bool        `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
}

type OAuthResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Data    string `json:"data"`
	Code    int    `json:"code"`
}

type LeaveFile struct {
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	Name       string `json:"name"`
	ID         string `json:"id"`
	StorageKey string `json:"storageKey"`
	URL        string `json:"url"`
	OAID       int64  `json:"oaId"`
}

type LeaveSaveData struct {
	ClassCode    string                 `json:"classCode"`
	Hours        float64                `json:"hours"`
	LeaveModel   int                    `json:"leaveModel"`
	EndDate      string                 `json:"endDate"`
	CustomFields map[string]interface{} `json:"customFields"`
	EmployeeID   string                 `json:"employeeId"`
	Times        int                    `json:"times"`
	LeaveReason  string                 `json:"leaveReason"`
	Repeat       int                    `json:"repeat"`
	Files        []LeaveFile            `json:"files"`
	StartTime    string                 `json:"startTime"`
	EndTime      string                 `json:"endTime"`
	ReasonCode   string                 `json:"reasonCode"`
	Detail       string                 `json:"detail"`
	AutoCalc     bool                   `json:"autoCalc"`
	StartDate    string                 `json:"startDate"`
}

type LeaveSubmitRequest struct {
	ProcessInstanceID string          `json:"processInstanceId"`
	SubmitCheck       string          `json:"submitCheck"`
	MultipleApplyMode string          `json:"multipleApplyMode"`
	SaveData          []LeaveSaveData `json:"saveData"`
	ApplyEmployeeID   string          `json:"applyEmployeeId"`
	FormNo            string          `json:"formNo"`
}

type MyApplicationsRequest struct {
	StartDate  string `form:"startDate" json:"startDate"`
	EndDate    string `form:"endDate" json:"endDate"`
	FormType   string `form:"formType" json:"formType"`
	EmployeeID string `form:"employeeId" json:"employeeId"`
	FormNo     string `form:"formNo" json:"formNo"`
	Status     string `form:"status" json:"status"`
}

type LeaveQuotaRequest struct {
	Size             int    `json:"size"`
	EndDate          string `json:"endDate"`
	UnitCode         string `json:"unitCode"`
	EmployeeID       string `json:"employeeId"`
	Page             int    `json:"page"`
	IsIncludeSubUnit bool   `json:"isIncludeSubUnit"`
	StartDate        string `json:"startDate"`
}

type LeaveTypesRequest struct {
	EmployeeID      string `form:"employeeId" json:"employeeId" binding:"required"`
	Date            string `form:"date" json:"date" binding:"required"`
	BatchApplyLeave string `form:"batchApplyLeave" json:"batchApplyLeave" binding:"required"`
}

type LeaveHoursRequest struct {
	Data map[string]interface{} `json:"data"`
}

type ExceptionListRequest struct {
	StartDate  string `form:"startDate" json:"startDate"`
	EndDate    string `form:"endDate" json:"endDate"`
	EmployeeID string `form:"employeeId" json:"employeeId"`
}
