package assert

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func Equal[T comparable](t *testing.T, got, want T, context string) {
	t.Helper()

	if got != want {
		t.Errorf("%s: got: %v; want: %v", context, got, want)
	}
}

func NotEqual[T comparable](t *testing.T, got, want T, context string) {
	t.Helper()

	if got == want {
		t.Errorf("%s: got: %v; want: %v", context, got, want)
	}
}

func Unmarshal[T any](t *testing.T, rawPayload string, want *T) {
	t.Helper()

	if err := json.Unmarshal([]byte(rawPayload), want); err != nil {
		t.Fatalf("could not parse %q into %T: %v", rawPayload, want, err)
	}
}

func JSONDecode[T any](t *testing.T, b *bytes.Buffer, v *T) {
	t.Helper()

	if err := json.NewDecoder(b).Decode(v); err != nil {
		t.Fatalf("failed to decode JSON: %T in Response body: %v due to error: %v", v, b, err)
	}
}

func TimeAfter(t *testing.T, got, want time.Time, context string) {
	t.Helper()

	if !got.After(want) {
		t.Errorf("%s: got %v, want %v", context, got, want)
	}
}
