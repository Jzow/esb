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
	Size             *int   `json:"size,omitempty"`
	EndDate          string `json:"endDate,omitempty"`
	UnitCode         string `json:"unitCode,omitempty"`
	EmployeeID       string `json:"employeeId,omitempty"`
	Page             *int   `json:"page,omitempty"`
	IsIncludeSubUnit *bool  `json:"isIncludeSubUnit,omitempty"`
	StartDate        string `json:"startDate,omitempty"`
}

type LeaveTypesRequest struct {
	EmploeeID       string `form:"emploeeId" json:"emploeeId"`
	EmployeeID      string `form:"employeeId" json:"-"`
	Date            string `form:"date" json:"date"`
	BatchApplyLeave string `form:"batchApplyLeave" json:"batchApplyLeave"`
}

type LeaveHoursRequest struct {
	PacketID                string                 `json:"packetId,omitempty"`
	IsEndDateManualChange   *bool                  `json:"isEndDateManualChange,omitempty"`
	EndDateLeaveModel       *int                   `json:"endDateLeaveModel,omitempty"`
	EndDate                 string                 `json:"endDate,omitempty"`
	LeaveModel              *int                   `json:"leaveModel,omitempty"`
	ReplaceForm             string                 `json:"replaceForm,omitempty"`
	StartDTM                string                 `json:"startdtm,omitempty"`
	AdvanceNum              string                 `json:"advanceNum,omitempty"`
	BreastfeedingLeaveModel string                 `json:"breastfeedingLeaveModel,omitempty"`
	AttendanceParentLabel   string                 `json:"attendanceParentLabel,omitempty"`
	IsShowTotalHours        *bool                  `json:"isShowTotalHours,omitempty"`
	Times                   *int                   `json:"times,omitempty"`
	EndModel                string                 `json:"endModel,omitempty"`
	Repeat                  *int                   `json:"repeat,omitempty"`
	StartTime               string                 `json:"startTime,omitempty"`
	StartModel              string                 `json:"startModel,omitempty"`
	AutoCalc                *bool                  `json:"autoCalc,omitempty"`
	RepeatEndDate           string                 `json:"repeatEndDate,omitempty"`
	IsReplaceClass          *int                   `json:"isReplaceClass,omitempty"`
	DelayNum                string                 `json:"delayNum,omitempty"`
	EndDTM                  string                 `json:"enddtm,omitempty"`
	ExtendedField           string                 `json:"extendedField,omitempty"`
	LeaveReason             string                 `json:"leaveReason,omitempty"`
	StartDateLeaveModel     *int                   `json:"startDateLeaveModel,omitempty"`
	PersonIDs               []string               `json:"personIds,omitempty"`
	PersonID                string                 `json:"personId,omitempty"`
	Detail                  string                 `json:"detail,omitempty"`
	EndTime                 string                 `json:"endTime,omitempty"`
	IsDateManualChange      *bool                  `json:"isDateManualChange,omitempty"`
	PayCode                 string                 `json:"payCode,omitempty"`
	ClassType               string                 `json:"classType,omitempty"`
	StartDate               string                 `json:"startDate,omitempty"`
	Extra                   map[string]interface{} `json:"-"`
}

type ExceptionListRequest struct {
	StartDate  string `form:"startDate" json:"startDate"`
	EndDate    string `form:"endDate" json:"endDate"`
	EmployeeID string `form:"employeeId" json:"employeeId"`
}
