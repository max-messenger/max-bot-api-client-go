package maxbot

import "github.com/max-messenger/max-bot-api-client-go/schemes"

// Keyboard implements a builder for the inline keyboard.
type Keyboard struct {
	rows []*KeyboardRow
}

// AddRow adds a row to the inline keyboard.
func (k *Keyboard) AddRow() *KeyboardRow {
	kr := &KeyboardRow{}
	k.rows = append(k.rows, kr)

	return kr
}

// Build returns the keyboard.
func (k *Keyboard) Build() schemes.Keyboard {
	buttons := make([][]schemes.ButtonInterface, 0, len(k.rows))
	for _, r := range k.rows {
		buttons = append(buttons, r.Build())
	}

	return schemes.Keyboard{Buttons: buttons}
}

// KeyboardRow represents a button row.
type KeyboardRow struct {
	cols []schemes.ButtonInterface
}

// Build returns keyboard rows.
func (k *KeyboardRow) Build() []schemes.ButtonInterface {
	return k.cols
}

func (k *KeyboardRow) AddButton(b schemes.ButtonInterface) *KeyboardRow {
	k.cols = append(k.cols, b)

	return k
}

// AddLink button.
func (k *KeyboardRow) AddLink(text string, _ schemes.Intent, url string) *KeyboardRow {
	b := schemes.LinkButton{
		Url: url,
		Button: schemes.Button{
			Text: text,
			Type: schemes.LINK,
		},
	}
	k.cols = append(k.cols, b)

	return k
}

// AddCallback button.
func (k *KeyboardRow) AddCallback(text string, intent schemes.Intent, payload string) *KeyboardRow {
	b := schemes.CallbackButton{
		Payload: payload,
		Intent:  intent,
		Button: schemes.Button{
			Text: text,
			Type: schemes.CALLBACK,
		},
	}
	k.cols = append(k.cols, b)

	return k
}

// AddContact button.
func (k *KeyboardRow) AddContact(text string) *KeyboardRow {
	b := schemes.RequestContactButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.CONTACT,
		},
	}
	k.cols = append(k.cols, b)

	return k
}

// AddGeolocation button.
func (k *KeyboardRow) AddGeolocation(text string, quick bool) *KeyboardRow {
	b := schemes.RequestGeoLocationButton{
		Quick: quick,
		Button: schemes.Button{
			Text: text,
			Type: schemes.GEOLOCATION,
		},
	}
	k.cols = append(k.cols, b)

	return k
}

// AddOpenApp button.
func (k *KeyboardRow) AddOpenApp(text string, app, payload string, contactId int64) *KeyboardRow {
	b := schemes.OpenAppButton{
		WebApp:    app,
		Payload:   payload,
		ContactId: contactId,
		Button: schemes.Button{
			Text: text,
			Type: schemes.OPEN_APP,
		},
	}
	k.cols = append(k.cols, b)

	return k
}

// AddMessage button.
func (k *KeyboardRow) AddMessage(text string) *KeyboardRow {
	b := schemes.MessageButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.MESSAGE,
		},
	}
	k.cols = append(k.cols, b)

	return k
}

func (k *KeyboardRow) AddClipboard(text, payload string) *KeyboardRow {
	b := schemes.ClipboardButton{
		Payload: payload,
		Button: schemes.Button{
			Text: text,
			Type: schemes.CLIPBOARD,
		},
	}
	k.cols = append(k.cols, b)

	return k
}

// keyboards sugar

// InlineKeyboard create inline keyboard
func InlineKeyboard(rows ...[]schemes.ButtonInterface) *Keyboard {
	k := &Keyboard{
		rows: make([]*KeyboardRow, 0, len(rows)),
	}
	for _, r := range rows {
		k.rows = append(k.rows, &KeyboardRow{cols: r})
	}
	return k
}

// Row creates one line of buttons.
func Row(buttons ...schemes.ButtonInterface) []schemes.ButtonInterface {
	return buttons
}

// Btn creates a standard Callback button.
// Optionally, you can pass schemes.Intent (default is schemes.DEFAULT).
func Btn(text, payload string, intent ...schemes.Intent) schemes.ButtonInterface {
	i := schemes.DEFAULT
	if len(intent) > 0 {
		i = intent[0]
	}
	return schemes.CallbackButton{
		Payload: payload,
		Intent:  i,
		Button: schemes.Button{
			Text: text,
			Type: schemes.CALLBACK,
		},
	}
}

// BtnLink creates a link button.
func BtnLink(text, url string) schemes.ButtonInterface {
	return schemes.LinkButton{
		Url: url,
		Button: schemes.Button{
			Text: text,
			Type: schemes.LINK,
		},
	}
}

// BtnContact creates a contact request button.
func BtnContact(text string) schemes.ButtonInterface {
	return schemes.RequestContactButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.CONTACT,
		},
	}
}

// BtnGeo creates a location request button.
func BtnGeo(text string, quick bool) schemes.ButtonInterface {
	return schemes.RequestGeoLocationButton{
		Quick: quick,
		Button: schemes.Button{
			Text: text,
			Type: schemes.GEOLOCATION,
		},
	}
}

// BtnApp creates a button to open the Web App.
func BtnApp(text, app, payload string, contactId int64) schemes.ButtonInterface {
	return schemes.OpenAppButton{
		WebApp:    app,
		Payload:   payload,
		ContactId: contactId,
		Button: schemes.Button{
			Text: text,
			Type: schemes.OPEN_APP,
		},
	}
}

// BtnMsg creates a button to send text on behalf of the user.
func BtnMsg(text string) schemes.ButtonInterface {
	return schemes.MessageButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.MESSAGE,
		},
	}
}

// BtnClipboard creates a button to copy to the clipboard.
func BtnClipboard(text, payload string) schemes.ButtonInterface {
	return schemes.ClipboardButton{
		Payload: payload,
		Button: schemes.Button{
			Text: text,
			Type: schemes.CLIPBOARD,
		},
	}
}
