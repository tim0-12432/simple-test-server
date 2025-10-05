package docker

import (
	"testing"
	"time"
)

func TestParseFindOutput_Basic(t *testing.T) {
	raw := "f|123|1696425600.000000000|/usr/share/nginx/html/index.html\n" +
		"d|4096|1696425601.000000000|/usr/share/nginx/html/assets\n"
	entries, truncated := parseFindOutput(raw, "/usr/share/nginx/html", 100)
	if truncated {
		t.Fatalf("expected not truncated")
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Name != "index.html" || entries[0].Type != "file" || entries[0].Size != 123 {
		t.Fatalf("unexpected first entry: %+v", entries[0])
	}
	if !entries[0].ModifiedAt.Equal(time.Unix(1696425600, 0).UTC()) {
		t.Fatalf("unexpected modifiedAt: %v", entries[0].ModifiedAt)
	}
}

func TestParseFindOutput_MalformedAndTruncate(t *testing.T) {
	raw := "badlinewithoutpipes\n" +
		"f|10|1696425600.5|/usr/share/nginx/html/.hidden\n" +
		"f|20|1696425601.5|/usr/share/nginx/html/visible.txt\n"
	entries, truncated := parseFindOutput(raw, "/usr/share/nginx/html", 1)
	if !truncated {
		t.Fatalf("expected truncated true")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry due to limit, got %d", len(entries))
	}
	if entries[0].Name != ".hidden" {
		t.Fatalf("expected hidden file name parsed, got %s", entries[0].Name)
	}
}
