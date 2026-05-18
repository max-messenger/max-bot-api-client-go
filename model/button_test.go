package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewKeyboard(t *testing.T) {
	k := NewKeyboard()

	if k == nil {
		t.Fatal("NewKeyboard returned nil")
	}

	if k.rows == nil {
		t.Error("Keyboard rows is nil, should be initialized empty slice")
	}

	if len(k.rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(k.rows))
	}
}

func TestKeyboardAddRow(t *testing.T) {
	k := NewKeyboard()

	row1 := k.AddRow()
	if row1 == nil {
		t.Error("AddRow returned nil")
	}

	if len(k.rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(k.rows))
	}

	row2 := k.AddRow()
	if len(k.rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(k.rows))
	}

	// Verify rows are different
	if row1 == row2 {
		t.Error("AddRow returned same row instance")
	}
}

func TestKeyboardBuild(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Keyboard)
		expected [][]*Button
	}{
		{
			name:     "empty keyboard",
			setup:    func(k *Keyboard) {},
			expected: [][]*Button{},
		},
		{
			name: "single row no buttons",
			setup: func(k *Keyboard) {
				k.AddRow()
			},
			expected: [][]*Button{{}},
		},
		{
			name: "single row multiple buttons",
			setup: func(k *Keyboard) {
				row := k.AddRow()
				row.AddMessage("Btn1")
				row.AddMessage("Btn2")
			},
			expected: func() [][]*Button {
				return [][]*Button{
					{
						{Text: "Btn1", Type: ButtonMessage},
						{Text: "Btn2", Type: ButtonMessage},
					},
				}
			}(),
		},
		{
			name: "multiple rows",
			setup: func(k *Keyboard) {
				row1 := k.AddRow()
				row1.AddMessage("A")
				row2 := k.AddRow()
				row2.AddMessage("B")
				row2.AddMessage("C")
			},
			expected: func() [][]*Button {
				return [][]*Button{
					{{Text: "A", Type: ButtonMessage}},
					{
						{Text: "B", Type: ButtonMessage},
						{Text: "C", Type: ButtonMessage},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeyboard()
			tt.setup(k)

			result := k.Build()

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d rows, got %d", len(tt.expected), len(result))
				return
			}

			for i := 0; i < len(result); i++ {
				if len(result[i]) != len(tt.expected[i]) {
					t.Errorf("Row %d: expected %d buttons, got %d",
						i, len(tt.expected[i]), len(result[i]))
				}
			}
		})
	}
}

func TestKeyboardRowAddCallback(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		intent   Intent
		payload  string
		expected Button
	}{
		{
			name:    "basic callback",
			text:    "Click me",
			intent:  IntentPositive,
			payload: "action_1",
			expected: Button{
				Type:    ButtonCallback,
				Text:    "Click me",
				Intent:  IntentPositive,
				Payload: "action_1",
			},
		},
		{
			name:    "empty payload",
			text:    "No payload",
			intent:  IntentPositive,
			payload: "",
			expected: Button{
				Type:    ButtonCallback,
				Text:    "No payload",
				Intent:  IntentPositive,
				Payload: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &KeyboardRow{cols: make([]*Button, 0)}
			result := row.AddCallback(tt.text, tt.intent, tt.payload)

			if result != row {
				t.Error("AddCallback should return the same row instance for chaining")
			}

			if len(row.cols) != 1 {
				t.Fatalf("Expected 1 button, got %d", len(row.cols))
			}

			button := row.cols[0]
			if button.Type != tt.expected.Type {
				t.Errorf("Type: expected %v, got %v", tt.expected.Type, button.Type)
			}
			if button.Text != tt.expected.Text {
				t.Errorf("Text: expected %s, got %s", tt.expected.Text, button.Text)
			}
			if button.Intent != tt.expected.Intent {
				t.Errorf("Intent: expected %v, got %v", tt.expected.Intent, button.Intent)
			}
			if button.Payload != tt.expected.Payload {
				t.Errorf("Payload: expected %s, got %s", tt.expected.Payload, button.Payload)
			}
		})
	}
}

func TestKeyboardRowAddLink(t *testing.T) {
	row := &KeyboardRow{cols: make([]*Button, 0)}
	result := row.AddLink("Google", "https://google.com")

	if result != row {
		t.Error("AddLink should return the same row instance for chaining")
	}

	if len(row.cols) != 1 {
		t.Fatalf("Expected 1 button, got %d", len(row.cols))
	}

	button := row.cols[0]
	if button.Type != ButtonLink {
		t.Errorf("Type: expected %v, got %v", ButtonLink, button.Type)
	}
	if button.Text != "Google" {
		t.Errorf("Text: expected Google, got %s", button.Text)
	}
	if button.URL != "https://google.com" {
		t.Errorf("URL: expected https://google.com, got %s", button.URL)
	}
}

func TestKeyboardRowAddGeoLocation(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		quick bool
	}{
		{"quick geo", "Send Location", true},
		{"regular geo", "Share Location", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &KeyboardRow{cols: make([]*Button, 0)}
			result := row.AddGeoLocation(tt.text, tt.quick)

			if result != row {
				t.Error("AddGeoLocation should return the same row instance for chaining")
			}

			button := row.cols[0]
			if button.Type != ButtonRequestGeo {
				t.Errorf("Type: expected %v, got %v", ButtonRequestGeo, button.Type)
			}
			if button.Text != tt.text {
				t.Errorf("Text: expected %s, got %s", tt.text, button.Text)
			}
			if button.Quick != tt.quick {
				t.Errorf("Quick: expected %v, got %v", tt.quick, button.Quick)
			}
		})
	}
}

func TestKeyboardRowAddContact(t *testing.T) {
	row := &KeyboardRow{cols: make([]*Button, 0)}
	result := row.AddContact("Share Contact")

	if result != row {
		t.Error("AddContact should return the same row instance for chaining")
	}

	button := row.cols[0]
	if button.Type != ButtonRequestContact {
		t.Errorf("Type: expected %v, got %v", ButtonRequestContact, button.Type)
	}
	if button.Text != "Share Contact" {
		t.Errorf("Text: expected 'Share Contact', got %s", button.Text)
	}
}

func TestKeyboardRowAddMessage(t *testing.T) {
	row := &KeyboardRow{cols: make([]*Button, 0)}
	result := row.AddMessage("Simple Message")

	if result != row {
		t.Error("AddMessage should return the same row instance for chaining")
	}

	button := row.cols[0]
	if button.Type != ButtonMessage {
		t.Errorf("Type: expected %v, got %v", ButtonMessage, button.Type)
	}
	if button.Text != "Simple Message" {
		t.Errorf("Text: expected 'Simple Message', got %s", button.Text)
	}
}

func TestKeyboardRowAddOpenApp(t *testing.T) {
	row := &KeyboardRow{cols: make([]*Button, 0)}
	result := row.AddOpenApp("Open App", 12345)

	if result != row {
		t.Error("AddOpenApp should return the same row instance for chaining")
	}

	button := row.cols[0]
	if button.Type != ButtonOpenApp {
		t.Errorf("Type: expected %v, got %v", ButtonOpenApp, button.Type)
	}
	if button.Text != "Open App" {
		t.Errorf("Text: expected 'Open App', got %s", button.Text)
	}
	if button.ContactID != 12345 {
		t.Errorf("ContactID: expected 12345, got %d", button.ContactID)
	}
	// WebApp should be set to text (as per implementation)
	if button.WebApp != "Open App" {
		t.Errorf("WebApp: expected 'Open App', got %s", button.WebApp)
	}
}

func TestKeyboardRowAddClipboard(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		payload string
	}{
		{"simple clipboard", "Copy", "text to copy"},
		{"empty payload", "Copy Empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &KeyboardRow{cols: make([]*Button, 0)}
			result := row.AddClipboard(tt.text, tt.payload)

			if result != row {
				t.Error("AddClipboard should return the same row instance for chaining")
			}

			button := row.cols[0]
			if button.Type != ButtonClipboard {
				t.Errorf("Type: expected %v, got %v", ButtonClipboard, button.Type)
			}
			if button.Text != tt.text {
				t.Errorf("Text: expected %s, got %s", tt.text, button.Text)
			}
			if button.Payload != tt.payload {
				t.Errorf("Payload: expected %s, got %s", tt.payload, button.Payload)
			}
		})
	}
}

func TestKeyboardRowChaining(t *testing.T) {
	row := &KeyboardRow{cols: make([]*Button, 0)}

	row.
		AddMessage("Msg1").
		AddCallback("Callback", IntentPositive, "data").
		AddLink("Link", "https://example.com").
		AddGeoLocation("Geo", true).
		AddContact("Contact").
		AddOpenApp("App", 999).
		AddClipboard("Copy", "secret")

	expectedCount := 7
	if len(row.cols) != expectedCount {
		t.Errorf("Expected %d buttons from chaining, got %d", expectedCount, len(row.cols))
	}

	// Verify types
	expectedTypes := []ButtonType{
		ButtonMessage,
		ButtonCallback,
		ButtonLink,
		ButtonRequestGeo,
		ButtonRequestContact,
		ButtonOpenApp,
		ButtonClipboard,
	}

	for i, expected := range expectedTypes {
		if row.cols[i].Type != expected {
			t.Errorf("Button %d: expected type %v, got %v", i, expected, row.cols[i].Type)
		}
	}
}

func TestKeyboardMultipleRowsAndButtons(t *testing.T) {
	k := NewKeyboard()

	// Row 1: 2 buttons
	row1 := k.AddRow()
	row1.AddMessage("Left").AddMessage("Right")

	// Row 2: 3 buttons
	row2 := k.AddRow()
	row2.AddLink("A", "a.com")
	row2.AddLink("B", "b.com")
	row2.AddLink("C", "c.com")

	// Row 3: 1 button
	row3 := k.AddRow()
	row3.AddCallback("Action", IntentPositive, "doIt")

	result := k.Build()

	if len(result) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(result))
	}

	if len(result[0]) != 2 {
		t.Errorf("Row1: expected 2 buttons, got %d", len(result[0]))
	}
	if len(result[1]) != 3 {
		t.Errorf("Row2: expected 3 buttons, got %d", len(result[1]))
	}
	if len(result[2]) != 1 {
		t.Errorf("Row3: expected 1 button, got %d", len(result[2]))
	}
}

func TestKeyboardRowBuild(t *testing.T) {
	row := &KeyboardRow{cols: make([]*Button, 0)}
	row.AddMessage("Test1").AddMessage("Test2")

	result := row.Build()

	if len(result) != 2 {
		t.Errorf("Expected 2 buttons, got %d", len(result))
	}

	// Verify the buttons are the same instances
	if result[0] != row.cols[0] {
		t.Error("Build should return the same button instances")
	}

	if result[1] != row.cols[1] {
		t.Error("Build should return the same button instances")
	}
}

// Optional: JSON marshaling test (if ButtonType implements MarshalJSON)
func TestButtonJSONMarshaling(t *testing.T) {
	btn := &Button{
		Text:    "Test",
		Type:    ButtonCallback,
		Intent:  IntentPositive,
		Payload: "test_payload",
	}

	data, err := json.Marshal(btn)
	if err != nil {
		t.Fatalf("Failed to marshal button: %v", err)
	}

	var result Button
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal button: %v", err)
	}

	if !reflect.DeepEqual(btn, &result) {
		t.Errorf("Round-trip mismatch.\nOriginal: %+v\nUnmarshaled: %+v", btn, &result)
	}
}
