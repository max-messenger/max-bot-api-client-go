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
