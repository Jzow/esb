package service_auth

import "testing"

func TestParseAppTokenTokenColonApps(t *testing.T) {
	token, apps := parseAppToken("token-a:gaiastandard|robot")
	if token != "token-a" {
		t.Fatalf("token = %q, want token-a", token)
	}
	if !allowApp(apps, "gaiastandard") || !allowApp(apps, "robot") {
		t.Fatalf("apps = %#v, want gaiastandard and robot", apps)
	}
	if allowApp(apps, "unknown") {
		t.Fatalf("apps = %#v should not allow unknown", apps)
	}
}

func TestParseAppTokenAppEqualsToken(t *testing.T) {
	token, apps := parseAppToken("gaiastandard=token-a")
	if token != "token-a" {
		t.Fatalf("token = %q, want token-a", token)
	}
	if !allowApp(apps, "gaiastandard") {
		t.Fatalf("apps = %#v, want gaiastandard", apps)
	}
}

func TestParseAppTokenPlainAllowsAll(t *testing.T) {
	token, apps := parseAppToken("token-a")
	if token != "token-a" {
		t.Fatalf("token = %q, want token-a", token)
	}
	if !allowApp(apps, "anything") {
		t.Fatalf("plain token should allow all apps")
	}
}
