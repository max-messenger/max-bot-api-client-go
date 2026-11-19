package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPayloadContactSerialization(t *testing.T) {
	data, err := json.Marshal(Payload{
		ContactID: 123456789,
		Name:      "John Doe",
		VCFInfo:   "BEGIN:VCARD",
		VCFPhone:  "+79991234567",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	wants := []string{
		`"contact_id":123456789`,
		`"name":"John Doe"`,
		`"vcf_info":"BEGIN:VCARD"`,
		`"vcf_phone":"+79991234567"`,
	}
	for _, want := range wants {
		if !strings.Contains(string(data), want) {
			t.Errorf("want %s in %s", want, data)
		}
	}
}

func TestPayloadPhotosSerialization(t *testing.T) {
	data, err := json.Marshal(Payload{
		Photos: map[string]PhotoToken{
			"0": {Token: "token-0"},
		},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	want := `{"photos":{"0":{"token":"token-0"}}}`
	if string(data) != want {
		t.Errorf("want %s, got %s", want, data)
	}
}

func TestPayloadEmptySerialization(t *testing.T) {
	data, err := json.Marshal(Payload{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(data) != "{}" {
		t.Errorf("want {}, got %s", data)
	}
}
