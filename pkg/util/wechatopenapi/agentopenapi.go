package wechatopenapi

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

const openApi = "https://qyapi.weixin.qq.com/cgi-bin"

type WechatAccessToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

//refer to https://developer.work.weixin.qq.com/document/path/91039
func GetAccessToken(corpid, corpsecret string) (*WechatAccessToken, []byte, error) {
	url := fmt.Sprintf("%s/gettoken?corpid=%s&corpsecret=%s", openApi, corpid, corpsecret)
	resp := new(WechatAccessToken)
	respBytes, err := util.Get(url, resp)
	/*	{
		"errcode": 0,
		"errmsg": "ok",
		"access_token": "EIVu2fRmpulS59eusDnRQEe-hXYy3cvhlvts9IxTyIY5ugpp3jrJIYWAL07cS67z6iZoHHcOURCQrsytKCFdY_z8fv9BIPomh2Z4oJ7S8fWfwXlnXtXqisUE0TNFRhusPSrr2PSHJUNyyXhWSPaHnXYnYkYd7GRDIucjPz4K76BD_-8f-NAP4j-g4369qN-3PN4CzvtYFBg8LLZhJ_qRsA",
		"expires_in": 7200
	}	*/
	return resp, respBytes, err
}

type WechatSyncMsgReq struct {
	Cursor      string `json:"cursor"`
	Token       string `json:"token"`
	Limit       int    `json:"limit"`
	VoiceFormat int    `json:"voice_format"`
	OpenKfid    string `json:"open_kfid"`
}
type WechatSyncMsgResp struct {
	ErrCode    int                      `json:"errcode"`
	ErrMsg     string                   `json:"errmsg"`
	NextCursor string                   `json:"next_cursor"`
	HasMore    int                      `json:"has_more"`
	MsgList    []map[string]interface{} `json:"msg_list"` //20多种消息类型，详情参考官方文档
}

//refer to https://developer.work.weixin.qq.com/document/path/94670
func SyncMsg(accessToken string, req WechatSyncMsgReq) (*WechatSyncMsgResp, []byte, error) {
	url := fmt.Sprintf("%s/kf/sync_msg?access_token=%s", openApi, accessToken)
	resp := new(WechatSyncMsgResp)
	respBytes, err := util.Post(url, req, resp)
	/*
	   {
	       "errcode": 0,
	       "errmsg": "ok",
	       "next_cursor": "4gw7MepFLfgF2VC5npv",
	       "msg_list": [
	           {
	               "msgid": "7ktYupMzEb6nVGvVogQKGP45Hz72kojQKkGmAUE16Ztq",
	               "send_time": 1685344525,
	               "origin": 4,
	               "msgtype": "event",
	               "event": {
	                   "event_type": "servicer_status_change",
	                   "servicer_userid": "jamky",
	                   "status": 1,
	                   "open_kfid": "wkvUrkEQAA-up5EHls8HdO1pDKumDFwA"
	               }
	           },
	           {
	               "msgid": "BHwpM4yofY28sgTrTFPpuBAW96",
	               "open_kfid": "wkvUrkEQAA-up5EHls8HdO1pDKumDFwA",
	               "external_userid": "wmvUrkEQAA8RYMrpStNlnFcMcImf35KQ",
	               "send_time": 1685608528,
	               "origin": 3,
	               "msgtype": "text",
	               "text": {
	                   "content": "测试问1"
	               }
	           }
	       ],
	       "has_more": 0
	   }*/
	return resp, respBytes, err
}

type WechatSendResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   string `json:"msgid"`
}

//refer to https://developer.work.weixin.qq.com/document/path/94677
func Send(accessToken string, req interface{}) (*WechatSendResp, []byte, error) {
	url := fmt.Sprintf("%s/kf/send_msg?access_token=%s", openApi, accessToken)
	resp := new(WechatSendResp)
	respBytes, err := util.Post(url, req, resp)
	/*
		{
		    "errcode": 0,
		    "errmsg": "ok",
		    "msgid": "BHwpM4yofY28sgTrTFPpuBAW96"
		}
		   }*/
	return resp, respBytes, err
}

func GetMsg(accessToken, openKfid, msgToken, logTag string, sendtime uint32) map[string]interface{} {
	hasData := true
	syncMsgReq := WechatSyncMsgReq{Token: msgToken, OpenKfid: openKfid}
	for hasData {
		resp, respBytes, err := SyncMsg(accessToken, syncMsgReq)
		if err != nil || resp == nil || resp.ErrCode != 0 {
			logging.Errorf("%s failed to wechatopenapi sync the message,result:%s", logTag, respBytes)
			return nil
		}

		//find data
		for _, msg := range resp.MsgList {
			if msg["msgtype"] == "text" && openKfid == util.ObjToStr(msg["open_kfid"]) && sendtime == util.ObjToUint32(msg["send_time"]) {
				return msg
			}
		}

		hasData = resp.HasMore == 1
		if hasData {
			syncMsgReq.Cursor = resp.NextCursor
		}
	}
	return nil
}
