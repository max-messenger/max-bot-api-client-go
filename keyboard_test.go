package maxbot

import (
	"reflect"
	"testing"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func TestKeyboardSugarEquivalence(t *testing.T) {
	// Original keyboard
	builderKb := &Keyboard{}

	row1 := builderKb.AddRow()
	row1.AddCallback("Yes", schemes.POSITIVE, "agree")
	row1.AddCallback("Nope", schemes.NEGATIVE, "decline")

	row2 := builderKb.AddRow()
	row2.AddLink("Site", schemes.DEFAULT, "https://example.com")
	row2.AddGeolocation("Location", true)

	// Sugar keyboard
	sugarKb := InlineKeyboard(
		Row(
			Btn("Yes", "agree", schemes.POSITIVE),
			Btn("Nope", "decline", schemes.NEGATIVE),
		),
		Row(
			BtnLink("Site", "https://example.com"),
			BtnGeo("Location", true),
		),
	)

	builderResult := builderKb.Build()
	sugarResult := sugarKb.Build()

	if !reflect.DeepEqual(builderResult, sugarResult) {
		t.Errorf("Keyboard structures mismatch!\nBuilder output: %+v\nSugar output:   %+v", builderResult, sugarResult)
	}
}
