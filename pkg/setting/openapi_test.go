package setting

import "testing"

func TestFindOpenAPIByURLMatchesConfiguredHost(t *testing.T) {
	old := OpenAPISettings
	defer func() { OpenAPISettings = old }()

	OpenAPISettings = map[string]*OpenAPI{
		"gaiastandard": {
			Name:    "gaiastandard",
			BaseUrl: "https://gaiaopenapi-s.copm.com.cn",
		},
	}

	openAPI, ok, err := FindOpenAPIByURL("/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list")
	if err != nil {
		t.Fatalf("FindOpenAPIByURL() error = %v", err)
	}
	if !ok {
		t.Fatal("FindOpenAPIByURL() ok = false, want true")
	}
	if openAPI.Name != "gaiastandard" {
		t.Fatalf("openAPI.Name = %q, want gaiastandard", openAPI.Name)
	}
}

func TestFindOpenAPIByURLRejectsUnknownHost(t *testing.T) {
	old := OpenAPISettings
	defer func() { OpenAPISettings = old }()

	OpenAPISettings = map[string]*OpenAPI{
		"gaiastandard": {
			Name:    "gaiastandard",
			BaseUrl: "https://gaiaopenapi-s.copm.com.cn",
		},
	}

	_, ok, err := FindOpenAPIByURL("/https://example.com/anything")
	if err != nil {
		t.Fatalf("FindOpenAPIByURL() error = %v", err)
	}
	if ok {
		t.Fatal("FindOpenAPIByURL() ok = true, want false")
	}
}
