package godiffy

import (
	"reflect"
	"testing"
)

func TestParseSingleHunk(t *testing.T) {
	input := `diff --git a/foo.txt b/foo.txt
index abc123..def456 100644
--- a/foo.txt
+++ b/foo.txt
@@ -1,3 +1,4 @@
 line1
-line2
+new2
 line3
`

	diff, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(diff.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(diff.Files))
	}

	file := diff.Files[0]
	if file.OldPath != "foo.txt" {
		t.Errorf("expected OldPath 'foo.txt', got '%s'", file.OldPath)
	}
	if file.NewPath != "foo.txt" {
		t.Errorf("expected NewPath 'foo.txt', got '%s'", file.NewPath)
	}

	if len(file.Hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(file.Hunks))
	}

	h := file.Hunks[0]
	if h.OldStart != 1 {
		t.Errorf("expected OldStart 1, got %d", h.OldStart)
	}
	if h.OldLineCount != 3 {
		t.Errorf("expected OldLineCount 3, got %d", h.OldLineCount)
	}
	if h.NewStart != 1 {
		t.Errorf("expected NewStart 1, got %d", h.NewStart)
	}
	if h.NewLineCount != 4 {
		t.Errorf("expected NewLineCount 4, got %d", h.NewLineCount)
	}

	expectedLines := []*HunkLine{
		{Type: HunkLineContext, Content: "line1\n"},
		{Type: HunkLineDeleted, Content: "line2\n"},
		{Type: HunkLineAdded, Content: "new2\n"},
		{Type: HunkLineContext, Content: "line3\n"},
	}
	if !reflect.DeepEqual(h.Lines, expectedLines) {
		t.Errorf("h.Lines = %+v, want %+v", h.Lines, expectedLines)
	}
}

func TestParseMultipleHunks(t *testing.T) {
	input := `diff --git a/foo.txt b/foo.txt
index abc123..def456 100644
--- a/foo.txt
+++ b/foo.txt
@@ -1,2 +1,2 @@
-line1
+ONE
@@ -4,2 +5,3 @@
 line4
+four point one
 line5
`

	diff, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(diff.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(diff.Files))
	}

	file := diff.Files[0]
	if len(file.Hunks) != 2 {
		t.Fatalf("expected 2 hunks, got %d", len(file.Hunks))
	}

	// Verify first hunk
	h1 := file.Hunks[0]
	if h1.OldStart != 1 || h1.OldLineCount != 2 {
		t.Errorf("h1 start/count = %d/%d, want 1/2", h1.OldStart, h1.OldLineCount)
	}
	if h1.NewStart != 1 || h1.NewLineCount != 2 {
		t.Errorf("h1 start/count = %d/%d, want 1/2", h1.NewStart, h1.NewLineCount)
	}

	// Verify second hunk
	h2 := file.Hunks[1]
	if h2.OldStart != 4 || h2.OldLineCount != 2 {
		t.Errorf("h2 start/count = %d/%d, want 4/2", h2.OldStart, h2.OldLineCount)
	}
	if h2.NewStart != 5 || h2.NewLineCount != 3 {
		t.Errorf("h2 start/count = %d/%d, want 5/3", h2.NewStart, h2.NewLineCount)
	}
}
