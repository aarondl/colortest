package main

import "testing"

func TestChatty(t *testing.T) {
	t.Parallel()

	t.Log("talking")
}

func TestFail(t *testing.T) {
	t.Parallel()

	if 2 != 3 {
		t.Error("want 3")
	}
}

func TestPass(t *testing.T) {
	t.Parallel()
}

func TestSubTests(t *testing.T) {
	t.Parallel()

	t.Run("SubFail", func(t *testing.T) {
		t.Parallel()
		t.Error("failed")
	})
	t.Run("SubPass", func(t *testing.T) {
		t.Log("talking")
	})
}
