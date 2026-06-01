package openapi_proxy

import (
	"reflect"
	"testing"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

func TestBuildTargetURLRequiresFullURL(t *testing.T) {
	openAPI := &setting.OpenAPI{
		Name:       "gaiastandard",
		BaseUrl:    "https://example.com",
		CorpID:     "tenant-a",
		PathPrefix: "atd-appapi/api/v1/gaiastandard",
	}

	_, err := buildTargetURL(openAPI, "/getemployeeleaveremaindata/{tenant}", "page=1")
	if err == nil {
		t.Fatal("buildTargetURL() error = nil, want full url requirement")
	}
}

func TestBuildTargetURLKeepsAbsoluteURL(t *testing.T) {
	openAPI := &setting.OpenAPI{
		Name:       "gaiastandard",
		BaseUrl:    "https://gaiaopenapi-s.copm.com.cn",
		CorpID:     "zhwytest01",
		PathPrefix: "atd-appapi/api/v1/gaiastandard",
	}

	got, err := buildTargetURL(openAPI, "/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list", "employeeId=E00001")
	if err != nil {
		t.Fatalf("buildTargetURL() error = %v", err)
	}
	want := "https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list?employeeId=E00001"
	if got != want {
		t.Fatalf("buildTargetURL() = %q, want %q", got, want)
	}
}

func TestBuildTargetURLRejectsAbsoluteURLForDifferentHost(t *testing.T) {
	openAPI := &setting.OpenAPI{
		Name:    "gaiastandard",
		BaseUrl: "https://gaiaopenapi-s.copm.com.cn",
	}

	_, err := buildTargetURL(openAPI, "/https://example.com/anything", "")
	if err == nil {
		t.Fatal("buildTargetURL() error = nil, want different host rejection")
	}
}

func TestFixedHeadersReplacesCorpID(t *testing.T) {
	openAPI := &setting.OpenAPI{
		CorpID:       "zhwytest01",
		FixedHeaders: "tenant={CorpID},X-App=gaiastandard",
	}

	got := fixedHeaders(openAPI)
	if got["tenant"] != "zhwytest01" {
		t.Fatalf("tenant header = %q, want zhwytest01", got["tenant"])
	}
	if got["X-App"] != "gaiastandard" {
		t.Fatalf("X-App header = %q, want gaiastandard", got["X-App"])
	}
}

func TestOpenAPITokenKey(t *testing.T) {
	got := openAPITokenKey("gaiastandard")
	want := "esb:openapi:token:gaiastandard"
	if got != want {
		t.Fatalf("openAPITokenKey() = %q, want %q", got, want)
	}
}

func TestNormalizedTokenTTLDefaultsToTwoHours(t *testing.T) {
	got := normalizedTokenTTL(&setting.OpenAPI{})
	want := 2 * time.Hour
	if got != want {
		t.Fatalf("normalizedTokenTTL() = %s, want %s", got, want)
	}
}

func TestNormalizeResponseReturnsDetails(t *testing.T) {
	got, err := normalizeResponse([]byte(`{"code":200,"details":[{"id":"1"}],"message":"success"}`))
	if err != nil {
		t.Fatalf("normalizeResponse() error = %v", err)
	}
	want := []interface{}{map[string]interface{}{"id": "1"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeResponse() = %#v, want %#v", got, want)
	}
}

func TestNormalizeResponseReturnsData(t *testing.T) {
	got, err := normalizeResponse([]byte(`{"result":true,"data":{"total":1}}`))
	if err != nil {
		t.Fatalf("normalizeResponse() error = %v", err)
	}
	want := map[string]interface{}{"total": float64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeResponse() = %#v, want %#v", got, want)
	}
}
