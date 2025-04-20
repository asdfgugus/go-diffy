# GoDiffy
GoDiffy – Git diffs, parsed the Go way.

## Description
GoDiffy is a library that parses Git diffs and converts them into a structured format that can be easily consumed by Go programs. No external dependencies, just Go standard library.

## Example usage

```diff
diff --git a/foo.txt b/foo.txt
index abc123..def456 100644
--- a/foo.txt
+++ b/foo.txt
@@ -1,3 +1,4 @@
 line1
-line2
+new2
 line3
```

| Diff fragment                        | Parsed into               | Go field                                                                                   |
|--------------------------------------|----------------------------|--------------------------------------------------------------------------------------------|
| `diff --git a/foo.txt b/foo.txt`     | start of a new file block  | `FileDiff.Header`                                                                          |
| `index abc123..def456 100644`        | blob IDs + file mode       | `FileDiff.OldHash == "abc123"`<br>`FileDiff.NewHash == "def456"`<br>`FileDiff.NewMode == "100644"` |
| `--- a/foo.txt`                      | old filename               | `FileDiff.OldPath == "foo.txt"`                                                           |
| `+++ b/foo.txt`                      | new filename               | `FileDiff.NewPath == "foo.txt"`                                                           |
| `@@ -1,3 +1,4 @@`                    | hunk header:               | • old start & count: `Hunk.OldStart == 1`<br>  `Hunk.OldLineCount == 3`<br>• new start & count: `Hunk.NewStart == 1`<br>  `Hunk.NewLineCount == 4` |
| ` line1`                             | context (unchanged) line   | `HunkLine{Type: HunkLineContext, Line: "line1"}`                                           |
| `-line2`                             | deleted line               | `HunkLine{Type: HunkLineDeleted, Line: "line2"}`                                           |
| `+new2`                              | added line                 | `HunkLine{Type: HunkLineAdded, Line: "new2"}`                                              |
| ` line3`                             | context (unchanged) line   | `HunkLine{Type: HunkLineContext, Line: "line3"}`                                           |

Important
- Whitespace (` `, `-`, `+`) is stripped off before storing in HunkLine.Content.
- The order of hunks in FileDiff.Hunks matches the order of `@@ … @@` blocks.
- If you see multiple `@@ … @@` blocks, you’ll get multiple Hunk entries under the same FileDiff.
