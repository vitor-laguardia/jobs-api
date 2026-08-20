package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func assertEqual[T comparable](t *testing.T, got, want T, context string) {
	t.Helper()

	if got != want {
		t.Errorf("%s: got: %v; want: %v", context, got, want)
	}
}

func assertUnmarshal[T any](t *testing.T, rawPayload string, want *T) {
	t.Helper()

	if err := json.Unmarshal([]byte(rawPayload), want); err != nil {
		t.Fatalf("could not parse %q into %T: %v", rawPayload, want, err)
	}
}

func assertJSONDecode[T any](t *testing.T, b *bytes.Buffer, v *T) {
	t.Helper()

	if err := json.NewDecoder(b).Decode(v); err != nil {
		t.Fatalf("failed to decode JSON: %T in Response body: %v due to error: %v", v, b, err)
	}
}

func assertJob(t *testing.T, got Job, want Job) {
	t.Helper()

	if isEqual := reflect.DeepEqual(got, want); !isEqual {
		t.Errorf("jobs not equal: got: %#v, want: %#v", got, want)
	}
}
