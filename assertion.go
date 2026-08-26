package main

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func assertEqual[T comparable](t *testing.T, got, want T, context string) {
	t.Helper()

	if got != want {
		t.Errorf("%s: got: %v; want: %v", context, got, want)
	}
}

func assertNotEqual[T comparable](t *testing.T, got, want T, context string) {
	t.Helper()

	if got == want {
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

	if got.ID != want.ID ||
		got.Title != want.Title ||
		got.Description != want.Description ||
		got.Status != want.Status ||
		got.Priority != want.Priority ||
		got.UserID != want.UserID {
		t.Errorf("job fields mismatch: got %#v, want %#v", got, want)
	}

	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("createdAt mismatch: got %v, want %v", got.CreatedAt, want.CreatedAt)
	}
}

func assertTimeAfter(t *testing.T, got, want time.Time, context string) {
	t.Helper()

	if !got.After(want) {
		t.Errorf("%s: got %v, want %v", context, got, want)
	}
}
