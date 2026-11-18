package model

import (
	"encoding/json"
	"testing"
)

func TestSimpleQueryResultDecode(t *testing.T) {
	raw := `{"success":false,"message":"error message","code":"error.code"}`

	var res SimpleQueryResult
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if res.Success {
		t.Errorf("success want false, got true")
	}
	if res.Message != "error message" {
		t.Errorf("message want %q, got %q", "error message", res.Message)
	}
	if res.Code != "error.code" {
		t.Errorf("code want %q, got %q", "error.code", res.Code)
	}
}

func TestSimpleQueryResultEncodeWithoutCode(t *testing.T) {
	data, err := json.Marshal(SimpleQueryResult{Success: true})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(data) != `{"success":true}` {
		t.Errorf("want %q, got %q", `{"success":true}`, string(data))
	}
}
