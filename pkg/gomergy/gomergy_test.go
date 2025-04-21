package gomergy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asdfgugus/godiffy/pkg/godiffy"
)

//
// Integration tests for MergeToPath
//

func TestMergeToPath_ReadDirError(t *testing.T) {
	err := MergeToPath(&godiffy.Diff{}, "/path/does/not/exist")
	if err == nil {
		t.Fatal("expected error for unreadable directory, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read directory") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMergeToPath_DeleteAndSkip(t *testing.T) {
	dir := t.TempDir()
	toDelete := filepath.Join(dir, "foo.txt")
	if err := os.WriteFile(toDelete, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	diff := godiffy.Diff{Files: []*godiffy.FileDiff{
		{Status: godiffy.FileStatusDeleted, NewPath: "foo.txt"},
		{Status: godiffy.FileStatusDeleted, NewPath: "noexist.txt"},
	}}

	if err := MergeToPath(&diff, dir); err != nil {
		t.Fatalf("expected no error deleting existing & skipping non‑existent, got %v", err)
	}
	if _, err := os.Stat(toDelete); !os.IsNotExist(err) {
		t.Errorf("expected foo.txt to be removed, got %v", err)
	}
}

func TestMergeToPath_NewFile_WrapError(t *testing.T) {
	dir := t.TempDir()
	// conflict: make "conflict" a file so MkdirAll fails
	conflict := filepath.Join(dir, "conflict")
	if err := os.WriteFile(conflict, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	diff := godiffy.Diff{Files: []*godiffy.FileDiff{
		{
			Status:  godiffy.FileStatusNew,
			NewPath: "conflict/bar.txt",
			NewMode: "0644",
			Hunks: []*godiffy.Hunk{
				{Lines: []*godiffy.HunkLine{{Content: "x"}}},
			},
		},
	}}

	err := MergeToPath(&diff, dir)
	if err == nil {
		t.Fatal("expected mkdir error wrapped, got nil")
	}
	if !strings.Contains(err.Error(), "failed to handle new file conflict/bar.txt") {
		t.Errorf("wrong wrapper: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to create directory") {
		t.Errorf("wrong inner message: %v", err)
	}
}

func TestMergeToPath_NewFile_InvalidModeWrap(t *testing.T) {
	dir := t.TempDir()
	diff := godiffy.Diff{Files: []*godiffy.FileDiff{
		{
			Status:  godiffy.FileStatusNew,
			NewPath: "f.txt",
			NewMode: "nope",
			Hunks:   []*godiffy.Hunk{},
		},
	}}

	err := MergeToPath(&diff, dir)
	if err == nil {
		t.Fatal("expected mode‑parse error wrapped, got nil")
	}
	if !strings.Contains(err.Error(), "failed to handle new file f.txt") {
		t.Errorf("wrong wrapper: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to convert file mode") {
		t.Errorf("wrong inner message: %v", err)
	}
}

func TestMergeToPath_ModifyFile_Success(t *testing.T) {
	dir := t.TempDir()
	diff := godiffy.Diff{Files: []*godiffy.FileDiff{
		{
			Status:  godiffy.FileStatusModified,
			NewPath: "a/m.txt",
			NewMode: "0600",
			Hunks: []*godiffy.Hunk{
				{Lines: []*godiffy.HunkLine{
					{Type: godiffy.HunkLineDeleted, Content: "old\n"},
					{Type: godiffy.HunkLineContext, Content: "keep\n"},
					{Type: godiffy.HunkLineAdded, Content: "new\n"},
				}},
			},
		},
	}}

	if err := MergeToPath(&diff, dir); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	out := filepath.Join(dir, "a/m.txt")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != "keep\nnew\n" {
		t.Errorf("content = %q; want %q", got, "keep\nnew\n")
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if perms := info.Mode().Perm(); perms != 0600 {
		t.Errorf("perms = %v; want 0600", perms)
	}
}

func TestMergeToPath_ModifyFile_InvalidModeWrap(t *testing.T) {
	dir := t.TempDir()
	diff := godiffy.Diff{Files: []*godiffy.FileDiff{
		{
			Status:  godiffy.FileStatusModified,
			NewPath: "m.txt",
			NewMode: "bad",
			Hunks:   []*godiffy.Hunk{},
		},
	}}

	err := MergeToPath(&diff, dir)
	if err == nil {
		t.Fatal("expected mode‑parse error wrapped, got nil")
	}
	if !strings.Contains(err.Error(), "failed to handle modified file m.txt") {
		t.Errorf("wrong wrapper: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to convert file mode") {
		t.Errorf("wrong inner message: %v", err)
	}
}

func TestMergeToPath_ModifyFile_MkdirErrorWrap(t *testing.T) {
	dir := t.TempDir()
	// conflict path element
	if err := os.WriteFile(filepath.Join(dir, "foo"), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	diff := godiffy.Diff{Files: []*godiffy.FileDiff{
		{
			Status:  godiffy.FileStatusModified,
			NewPath: "foo/bar.txt",
			NewMode: "0644",
			Hunks: []*godiffy.Hunk{
				{Lines: []*godiffy.HunkLine{{Type: godiffy.HunkLineAdded, Content: "x"}}},
			},
		},
	}}

	err := MergeToPath(&diff, dir)
	if err == nil {
		t.Fatal("expected mkdir error wrapped, got nil")
	}
	if !strings.Contains(err.Error(), "failed to handle modified file foo/bar.txt") {
		t.Errorf("wrong wrapper: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to create directory") {
		t.Errorf("wrong inner message: %v", err)
	}
}

func TestMergeToPath_EmptyDiff(t *testing.T) {
	dir := t.TempDir()
	if err := MergeToPath(&godiffy.Diff{Files: nil}, dir); err != nil {
		t.Errorf("empty diff should succeed, got %v", err)
	}
}

//
// Unit tests for each handler
//

func TestHandleDeletedFile_DeleteAndSkip(t *testing.T) {
	dir := t.TempDir()
	f := &godiffy.FileDiff{NewPath: "x.txt"}
	fp := filepath.Join(dir, "x.txt")

	// skip non‑existent
	if err := handleDeletedFile(f, dir); err != nil {
		t.Fatal(err)
	}

	// create and delete
	if err := os.WriteFile(fp, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := handleDeletedFile(f, dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(fp); !os.IsNotExist(err) {
		t.Errorf("expected deleted, got %v", err)
	}
}

func TestHandleNewFile_SuccessAndInvalidMode(t *testing.T) {
	dir := t.TempDir()
	f := &godiffy.FileDiff{
		NewPath: "sub/z.txt",
		NewMode: "0644",
		Hunks: []*godiffy.Hunk{
			{Lines: []*godiffy.HunkLine{{Content: "hey"}}},
		},
	}

	// success
	if err := handleNewFile(f, dir); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	out := filepath.Join(dir, f.NewPath)
	data, _ := os.ReadFile(out)
	if string(data) != "hey" {
		t.Errorf("content = %q; want hey", data)
	}
	info, _ := os.Stat(out)
	if info.Mode().Perm() != 0644 {
		t.Errorf("perms = %v; want 0644", info.Mode().Perm())
	}

	// invalid mode
	f2 := &godiffy.FileDiff{
		NewPath: "sub/z.txt",
		NewMode: "bad",
		Hunks:   f.Hunks,
	}
	err := handleNewFile(f2, dir)
	if err == nil {
		t.Fatal("expected mode parse error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to convert file mode") {
		t.Errorf("wrong error: %v", err)
	}
}

func TestHandleModifiedFile_SuccessAndInvalidMode(t *testing.T) {
	dir := t.TempDir()
	f := &godiffy.FileDiff{
		NewPath: "q.txt",
		NewMode: "0600",
		Hunks: []*godiffy.Hunk{
			{Lines: []*godiffy.HunkLine{
				{Type: godiffy.HunkLineDeleted, Content: "d"},
				{Type: godiffy.HunkLineContext, Content: "c"},
				{Type: godiffy.HunkLineAdded, Content: "a"},
			}},
		},
	}

	// success
	if err := handleModifiedFile(f, dir); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	out := filepath.Join(dir, f.NewPath)
	data, _ := os.ReadFile(out)
	if string(data) != "ca" {
		t.Errorf("content = %q; want ca", data)
	}
	info, _ := os.Stat(out)
	if info.Mode().Perm() != 0600 {
		t.Errorf("perms = %v; want 0600", info.Mode().Perm())
	}

	// invalid mode
	f.NewMode = "oops"
	err := handleModifiedFile(f, dir)
	if err == nil {
		t.Fatal("expected mode parse error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to convert file mode") {
		t.Errorf("wrong error: %v", err)
	}
}
