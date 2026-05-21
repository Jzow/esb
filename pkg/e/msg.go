package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "invalid params",
	UNAUTHORIZED:   "unauthorized",
	FORBIDDEN:      "forbidden",

	ERROR_EXIST:       "item exist",
	ERROR_EXIST_FAIL:  "failed to check item exist",
	ERROR_NOT_EXIST:   "item not exist",
	ERROR_GET_FAIL:    "failed to get item",
	ERROR_COUNT_FAIL:  "failed to check item total",
	ERROR_ADD_FAIL:    "failed to add",
	ERROR_EDIT_FAIL:   "failed to edit",
	ERROR_DELETE_FAIL: "failed to delete",
	ERROR_EXPORT_FAIL: "failed to export",
	ERROR_IMPORT_FAIL: "failed to import",

	ERROR_AUTH_CHECK_FAIL:    "failed to check auth",
	ERROR_AUTH_CHECK_TIMEOUT: "authorize time out",
	ERROR_AUTH_GENERATE_FAIL: "failed to generate auth ticket",
	ERROR_AUTH_FAIL:          "failed to authorize",

	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "failed to save attachment",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "failed to check attachment path and auth",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "error attachment format",
}

// GetMsg get error information based on OperateCode
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
