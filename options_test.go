package maxbot

import "testing"

func TestDefaultErrorBufferSize(t *testing.T) {
	api, err := New("token")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if got := cap(api.client.errors); got != defaultErrorBufferSize {
		t.Fatalf("expected default error buffer %d, got %d", defaultErrorBufferSize, got)
	}
}

func TestWithErrorBufferSize(t *testing.T) {
	api, err := New("token", WithErrorBufferSize(64))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if got := cap(api.client.errors); got != 64 {
		t.Fatalf("expected error buffer 64, got %d", got)
	}
}

func TestWithErrorBufferSizeIgnoresNonPositive(t *testing.T) {
	api, err := New("token", WithErrorBufferSize(0))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if got := cap(api.client.errors); got != defaultErrorBufferSize {
		t.Fatalf("non-positive size must keep default %d, got %d", defaultErrorBufferSize, got)
	}
}

func TestNotifyErrorKeepsBurstWithinBuffer(t *testing.T) {
	api, err := New("token", WithErrorBufferSize(4))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	cl := api.client
	for i := 0; i < 4; i++ {
		cl.notifyError(errForTest(i))
	}
	// All four must be retained by the channel, not dropped to the log.
	for i := 0; i < 4; i++ {
		select {
		case <-api.GetErrors():
		default:
			t.Fatalf("expected %d buffered errors, channel drained early at %d", 4, i)
		}
	}
}

type errForTest int

func (e errForTest) Error() string { return "test error" }
