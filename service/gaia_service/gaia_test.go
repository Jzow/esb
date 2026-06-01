package gaia_service

import (
	"reflect"
	"testing"
)

func TestNormalizeResponseReturnsDetails(t *testing.T) {
	details := []interface{}{
		map[string]interface{}{"leaveCode": "A30"},
	}
	out := map[string]interface{}{
		"reason":  nil,
		"code":    float64(200),
		"details": details,
		"message": "success",
	}

	got, err := normalizeResponse(out)
	if err != nil {
		t.Fatalf("normalizeResponse() error = %v", err)
	}
	if !reflect.DeepEqual(got, details) {
		t.Fatalf("normalizeResponse() = %#v, want details", got)
	}
}

func TestNormalizeResponseReturnsData(t *testing.T) {
	data := map[string]interface{}{"total": float64(1)}
	out := map[string]interface{}{
		"result": true,
		"data":   data,
	}

	got, err := normalizeResponse(out)
	if err != nil {
		t.Fatalf("normalizeResponse() error = %v", err)
	}
	if !reflect.DeepEqual(got, data) {
		t.Fatalf("normalizeResponse() = %#v, want data", got)
	}
}

func TestNormalizeResponseKeepsBodyWithoutDetails(t *testing.T) {
	out := map[string]interface{}{
		"code":    float64(200),
		"message": "success",
	}

	got, err := normalizeResponse(out)
	if err != nil {
		t.Fatalf("normalizeResponse() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Fatalf("normalizeResponse() = %#v, want original body", got)
	}
}

func TestNormalizeResponseReturnsErrorForBusinessFailure(t *testing.T) {
	out := map[string]interface{}{
		"code":    float64(500),
		"message": "failed",
		"reason":  "bad request",
	}

	if _, err := normalizeResponse(out); err == nil {
		t.Fatal("normalizeResponse() error = nil, want business failure")
	}
}
