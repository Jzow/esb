# Gaia Proxy APIs

All APIs below are ESB service APIs. The service obtains Gaia OAuth tokens internally and then calls Gaia OpenAPI.

Base URL:

```text
http://127.0.0.1:8000
```

Authentication:

```http
Authorization: Bearer <ESB user token or app token>
```

For local testing, the default app token in `conf/app.ini` is:

```http
Authorization: Bearer esb-app-token-change-me
```

Standard ESB response:

```json
{
  "code": 200,
  "msg": "ok",
  "data": {}
}
```

For Gaia responses that contain `details`, ESB returns `details` as `data`.
For Gaia responses that contain `data`, ESB returns Gaia `data` as ESB `data`.

## Get Leave Types

Gaia API: `GET /atd-webapi/api/v1/attendance-type/leave/list`

ESB API:

```http
GET /api/gaia/v1/leave/types?employeeId=test&date=2026-05-10&batchApplyLeave=0
```

Query parameters:

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| employeeId | string | no | Employee ID. ESB forwards it to Gaia as `emploeeId` to match the Gaia document. |
| emploeeId | string | no | Gaia document spelling. Takes priority over `employeeId`. |
| date | string | no | Date, for example `2026-05-10`. |
| batchApplyLeave | string | no | Batch apply flag, for example `0`. |

Success `data` is an array of leave type objects.

## Get Leave Quota

Gaia API: `POST /atd-appapi/api/v1/gaiastandard/getemployeeleaveremaindata/{tenant}`

`{tenant}` is replaced by `CorpID` from `conf/app.ini`.

ESB API:

```http
POST /api/gaia/v1/leave/quota
Content-Type: application/json

{
  "size": 10,
  "endDate": "2024-05-30",
  "unitCode": "U001",
  "employeeId": "employee01",
  "page": 1,
  "isIncludeSubUnit": false,
  "startDate": "2024-05-30"
}
```

Success `data` is the Gaia `data` object:

```json
{
  "total": 1,
  "employeeDate": []
}
```

## Get Leave Hours

Gaia API: `POST /atd-webapi/api/v1/workflow/leave/hours`

ESB API:

```http
POST /api/gaia/v1/leave/hours
Content-Type: application/json

{
  "packetId": "uuid",
  "isEndDateManualChange": true,
  "endDateLeaveModel": 1,
  "endDate": "2026-05-05",
  "leaveModel": 1,
  "startTime": "08:00",
  "endTime": "18:00",
  "autoCalc": true,
  "personIds": ["1"],
  "personId": "xxx",
  "detail": "description",
  "payCode": "xxx",
  "classType": "xxx",
  "startDate": "2026-05-05"
}
```

Success `data` is the Gaia `details` object.

## Submit Leave

Gaia API: `POST /atd-webapi/api/v1/gaiastandard/openapi/apply-leave-list/leaveapplysubmit`

ESB API:

```http
POST /api/gaia/v1/leave/submit
Content-Type: application/json

{
  "processInstanceId": "xxx",
  "submitCheck": "2",
  "multipleApplyMode": "2",
  "saveData": [
    {
      "classCode": "AnnualLeave",
      "hours": 8,
      "leaveModel": 0,
      "endDate": "2026-05-20",
      "customFields": {},
      "employeeId": "E00002",
      "times": 1,
      "leaveReason": "reasonId",
      "repeat": 1,
      "files": [
        {
          "fileName": "attachment.pdf",
          "fileSize": 1024,
          "name": "attachment.pdf",
          "id": "fileId",
          "storageKey": "storageKey",
          "url": "https://example.com/file.pdf",
          "oaId": 123
        }
      ],
      "startTime": "09:00",
      "endTime": "18:00",
      "reasonCode": "ANNUAL",
      "detail": "leave reason",
      "autoCalc": false,
      "startDate": "2026-05-20"
    }
  ],
  "applyEmployeeId": "E00001",
  "formNo": "202605200001"
}
```

Success `data` is the Gaia `details` object:

```json
{
  "processInstanceId": "processInstanceId",
  "formNumber": "LEAVE202605200001",
  "warnMsg": "",
  "canSubmit": ""
}
```

## My Applications

Gaia API: `GET /atd-webapi/api/v1/workflow/apply/list`

ESB API:

```http
GET /api/gaia/v1/leave/my-applications?startDate=2026-05-01&endDate=2026-05-31&formType=leave&employeeId=E00001&formNo=LEAVE202605200001&status=APPROVING
```

Query parameters:

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| startDate | string | no | Start date. |
| endDate | string | no | End date. |
| formType | string | no | Form type, for example `leave`. |
| employeeId | string | no | Employee ID. |
| formNo | string | no | Form number. |
| status | string | no | Status, for example `APPROVING`. |

Success `data` is an array of application objects.

## Error Notes

If Gaia returns HTTP errors, ESB returns HTTP 500 with Gaia error text in `msg`.

If Gaia returns HTTP 200 but business fields indicate failure, for example `code != 200` or `result=false`, ESB also returns HTTP 500 with the Gaia message and reason.
