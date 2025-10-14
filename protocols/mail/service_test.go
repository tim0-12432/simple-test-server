package mail

import (
	"github.com/emersion/go-imap"
	"testing"
)

func TestDecodeHeader(t *testing.T) {
	s := "=?UTF-8?B?VGVzdCBTdWJqZWN0?="
	got := decodeHeader(s)
	if got == s || got == "" {
		t.Fatalf("expected decoded header, got %q", got)
	}
}

func TestFormatAddress(t *testing.T) {
	addr := &imap.Address{MailboxName: "john", HostName: "example.com", PersonalName: "John Doe"}
	got := formatAddress(addr)
	if got != "John Doe <john@example.com>" {
		t.Fatalf("unexpected formatted address: %s", got)
	}
}
