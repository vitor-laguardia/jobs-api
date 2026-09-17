package assert

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
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

func TimeEqual(t *testing.T, got, want time.Time, context string) {
	t.Helper()

	// TODO: truncate in user domain
	if !got.Truncate(time.Microsecond).Equal(want.Truncate(time.Microsecond)) {
		t.Errorf("%s, got: %v, want %v", context, got, want)
	}
}

func Nil(t *testing.T, got any) {
	t.Helper()
	if !isNil(got) {
		t.Errorf("got: %v; want: nil", got)
	}
}

func ErrorIs(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got: %v; want: %v", got, want)
	}
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	}
	return false
}
