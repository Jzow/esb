package userCenter_service

import (
	"context"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

type UserCenterTenantResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *struct {
		TenantId   string `json:"tenant_id"`
		TenantCode string `json:"tenant_code"`
		TenantName string `json:"tenant_name"`
		TenantType string `json:"tenant_type"`
	} `json:"data,omitempty"`
}

type UserCenterProjectResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *struct {
		ProjectId   string `json:"project_id"`
		ProjectCode string `json:"project_code"`
		ProjectName string `json:"project_name"`
		TenantId    string `json:"tenant_id"`
	} `json:"data,omitempty"`
}

type UserCenterTenantUserResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		TenantUserId     int64  `json:"tenant_user_id"`
		TenantId         string `json:"tenant_id"`
		UserId           string `json:"user_id"`
		TenantUserName   string `json:"tenant_user_name"`
		TenantUserCnName string `json:"tenant_user_cn_name"`
		Avatar           string `json:"avatar"`
		Telphone         string `json:"telphone"`
		Email            string `json:"email"`
	} `json:"data,omitempty"`
}

type UserCenterDeptUserInfoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *struct {
		TenantUserId     int64  `json:"tenant_user_id"`
		TenantId         string `json:"tenant_id"`
		UserId           string `json:"user_id"`
		TenantUserName   string `json:"tenant_user_name"`
		TenantUserCnName string `json:"tenant_user_cn_name"`
		Avatar           string `json:"avatar"`
		Branch           string `json:"branch"`
		Telphone         string `json:"telphone"`
		Email            string `json:"email"`
		Status           int    `json:"status"`
		UserName         string `json:"user_name"`
		UserCnName       string `json:"user_cn_name"`
		WxUnionId        string `json:"wx_union_id"`

		Dept []struct {
			DeptId   int64  `json:"dept_id"`
			DeptName string `json:"dept_name"`
			IsMain   int    `json:"is_main"`
		} `json:"dept"`
	} `json:"data,omitempty"`
}

type UserCenterDeptResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *struct {
		DeptId   string
		ParentId string
		DeptName string
		TenantId string
	} `json:"data,omitempty"`
}

func GetTenantInfo(ctx context.Context, tenantId string) *UserCenterTenantResp {
	url := fmt.Sprintf("%v/api/system/v1/tenant/%v", setting.AppSetting.UserCenterApiHost, tenantId)
	resp := new(UserCenterTenantResp)
	respBytes, err := util.Get(url, util.GetServerAuthorization(), resp)
	util.Log("[GetTenantInfo]url=%v,result:%v", ctx, url, respBytes, err)
	if err == nil && resp != nil && resp.Code == 200 && resp.Data != nil {
		return resp
	}
	return nil
}

func GetProjectInfo(projectId string) *UserCenterProjectResp {
	url := fmt.Sprintf("%s/api/system/v1/project/%s", setting.AppSetting.UserCenterApiHost, projectId)
	resp := new(UserCenterProjectResp)
	respBytes, err := util.Get(url, util.GetServerAuthorization(), resp)
	logging.Infof("[GetProjectInfo]url=%s,resp:%s,error:%v", url, respBytes, err)
	if err == nil && resp != nil && resp.Code == 200 && resp.Data != nil {
		return resp
	}
	return nil
}

func GetTenantUserInfo(tenantId, userId string) *UserCenterTenantUserResp {
	url := fmt.Sprintf("%s/api/system/v1/tenantuser?tenant_id=%s&user_id=%s", setting.AppSetting.UserCenterApiHost, tenantId, userId)
	resp := new(UserCenterTenantUserResp)
	respBytes, err := util.Get(url, util.GetServerAuthorization(), resp)
	logging.Infof("[GetTenantUserInfo]url=%s,resp:%s,error:%v", url, respBytes, err)
	if err == nil && resp != nil && resp.Code == 200 && resp.Data != nil && len(resp.Data) > 0 {
		return resp
	}
	return nil
}

func GetTenantUsers(Token, tenantId, userId string) *UserCenterTenantUserResp {
	url := fmt.Sprintf("%s/api/system/v1/tenantuser?tenant_id=%s&tenant_user_id=%s", setting.AppSetting.UserCenterApiHost, tenantId, userId)
	resp := new(UserCenterTenantUserResp)
	header := map[string]any{"Authorization": Token, "X-Tenant-Id": tenantId}
	respBytes, err := util.Get(url, header, resp)
	logging.Infof("[GetTenantUserInfo]url=%s,resp:%s,error:%v", url, respBytes, err)
	if err == nil && resp != nil && resp.Code == 200 && resp.Data != nil && len(resp.Data) > 0 {
		return resp
	}
	return nil
}
func GetTenantUserInfoWithDept(tenantId, userId string) *UserCenterDeptUserInfoResp {
	url := fmt.Sprintf("%s/api/system/v1/tenant/user/info?tenant_id=%s&user_id=%s", setting.AppSetting.UserCenterApiHost, tenantId, userId)
	resp := new(UserCenterDeptUserInfoResp)
	respBytes, err := util.Get(url, util.GetServerAuthorization(), resp)
	logging.Infof("[GetTenantUserInfoWithDept]url=%s,resp:%s,error:%v", url, respBytes, err)
	if err == nil && resp != nil && resp.Code == 200 && resp.Data != nil {
		return resp
	}
	return nil
}

func GetTenantUserName(tenantId string, userId []string) map[string]string {
	url := fmt.Sprintf("%s/api/system/v1/users", setting.AppSetting.UserCenterApiHost)
	resp := new(models.ApiResp[[]map[string]any])
	respBytes, err := util.Post(url, map[string]any{"tenant_id": tenantId, "user_id": userId}, util.GetServerAuthorization(), resp)
	logging.Infof("[GetTenantUserName]url=%s,resp:%s,error:%v", url, respBytes, err)
	userIdNameMap := make(map[string]string)
	if err == nil && resp != nil && resp.Code == 200 && resp.Data != nil && len(resp.Data) > 0 {
		for _, itemMap := range resp.Data {
			if uid, ok := itemMap["user_id"]; ok {
				userIdNameMap[util.ObjToStr(uid)] = util.ObjToStr(itemMap["tenant_user_cn_name"])
			}
		}
	}
	return userIdNameMap
}

func GetTenantUserAvatar(tenantId string, userId []string) map[string]string {
	url := fmt.Sprintf("%s/api/system/v1/users", setting.AppSetting.UserCenterApiHost)
	resp := new(models.ApiResp[[]map[string]any])
	respBytes, err := util.Post(url, map[string]any{"tenant_id": tenantId, "user_id": userId}, util.GetServerAuthorization(), resp)
	logging.Infof("[GetTenantUserName]url=%s,resp:%s,error:%v", url, respBytes, err)
	userIdNameMap := make(map[string]string)
	if err == nil && resp != nil && resp.Code == 200 && resp.Data != nil && len(resp.Data) > 0 {
		for _, itemMap := range resp.Data {
			if uid, ok := itemMap["user_id"]; ok {
				userIdNameMap[util.ObjToStr(uid)] = util.ObjToStr(itemMap["avatar"])
			}
		}
	}
	return userIdNameMap
}

func GetRoleUsersByCode(ctx context.Context, tenantId, projectId, roleCode string) (userId []string) {
	url := fmt.Sprintf("%v/api/system/v1/rolecodeuser?role_code=%v&tenant_id=%v&project_id=%v",
		setting.AppSetting.UserCenterApiHost, roleCode, tenantId, projectId)
	resp := new(models.ApiResp[[]string])
	respBytes, err := util.Get(url, util.GetServerAuthorization(), resp)
	util.Log("[GetRoleUsersByCode]url=%v,result:%v", ctx, url, respBytes, err)
	if err == nil && resp != nil && resp.Code == 200 {
		userId = resp.Data
	}
	return
}

func CheckSignature(Token, tenantId, userId string) *CheckSignatureResp {
	url := fmt.Sprintf("%s/api/system/v1/check_signature?tenant_id=%s&user_id=%s&token=%s", setting.AppSetting.UserCenterApiHost, tenantId, userId, Token)
	resp := new(CheckSignatureResp)
	header := map[string]any{"Authorization": Token, "X-Tenant-Id": tenantId}
	respBytes, err := util.Get(url, header, resp)
	logging.Infof("[GetTenantUserInfo]url=%s,resp:%s,error:%v", url, respBytes, err)
	if err == nil && resp != nil && resp.Code == 200 {
		return resp
	}
	return nil
}

type CheckSignatureResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
