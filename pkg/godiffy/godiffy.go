package godiffy

import (
	"fmt"
	"strings"
)

func Parse(input string) (*Diff, error) {
	resultDiff := &Diff{}
	var currentFile *FileDiff

	for line := range strings.Lines(input) {
		if currentFile == nil && !strings.HasPrefix(line, "diff --git") {
			continue
		}
		switch {
		case strings.HasPrefix(line, "di"): // diff --git a/foo.txt b/foo.txt
			currentFile = parseNewFileDiff(resultDiff, line)
		case strings.HasPrefix(line, "i"): // index abc123..def456 100644
			parseMetadata(currentFile, line)
		case strings.HasPrefix(line, "--"): // --- a/foo.txt
			parseOldFilenameMarker(currentFile, line)
		case strings.HasPrefix(line, "++"): // +++ b/foo.txt
			parseNewFilenameMarker(currentFile, line)
		case strings.HasPrefix(line, "@"): // @@ -1,3 +1,4 @@
			parseHunk(currentFile, line)
		case strings.HasPrefix(line, "+"): // +added line
			parseAddLine(currentFile, line)
		case strings.HasPrefix(line, "-"): // -removed line
			parseDeleteLine(currentFile, line)
		case strings.HasPrefix(line, " "): //  context line
			parseContentLine(currentFile, line)
		case strings.HasPrefix(line, "new f"): // new file mode 100644
			parseNewFileMode(currentFile, line)
		case strings.HasPrefix(line, "de"): // deleted file mode 100644
			parseDeletedFileMode(currentFile, line)
		case strings.HasPrefix(line, "r"): // rename from old_name.txt / rename to new_name.txt
			if strings.HasPrefix(line, "rename f") {
				parseRenameFrom(currentFile, line)
			} else {
				parseRenameTo(currentFile, line)
			}
		case strings.HasPrefix(line, "o"): // old mode 100755
			parseOldMode(currentFile, line)
		case strings.HasPrefix(line, "new m"): // new mode 100644
			parseNewMode(currentFile, line)
		case strings.TrimSpace(line) == "": // (blank)
			continue
		default:
			return nil, fmt.Errorf("failed to parse line: %s", line)
		}
	}
	return resultDiff, nil
}

func parseNewFileDiff(resultDiff *Diff, line string) (currentFile *FileDiff) {
	currentFile = &FileDiff{
		Header: line,
	}
	resultDiff.Files = append(resultDiff.Files, currentFile)
	return currentFile

}

func parseHunk(currentFile *FileDiff, line string) {
	return
}

func parseOldFilenameMarker(currentFile *FileDiff, line string) {
	currentFile.OldPath = strings.SplitN(line, "/", 2)[1]
}

func parseNewFilenameMarker(currentFile *FileDiff, line string) {
	currentFile.NewPath = strings.SplitN(line, "/", 2)[1]
}

func parseAddLine(currentFile *FileDiff, line string) {
	return
}

func parseDeleteLine(currentFile *FileDiff, line string) {
	return
}

func parseContentLine(currentFile *FileDiff, line string) {
	return
}

func parseMetadata(currentFile *FileDiff, line string) {
	parts := strings.Fields(line) // ["index","abc123..def456","100644"]
	hashes := strings.SplitN(parts[1], "..", 2)
	currentFile.OldHash, currentFile.NewHash = hashes[0], hashes[1]
	currentFile.NewMode = parts[2]
}

func parseNewFileMode(currentFile *FileDiff, line string) {
	currentFile.IsNew = true
	currentFile.NewMode = strings.Fields(line)[3]
}

func parseDeletedFileMode(currentFile *FileDiff, line string) {
	currentFile.IsDeleted = true
	currentFile.OldMode = strings.Fields(line)[3]
}

func parseRenameFrom(currentFile *FileDiff, line string) {
	currentFile.IsRename = true
	currentFile.OldName = strings.Fields(line)[2]
}

func parseRenameTo(currentFile *FileDiff, line string) {
	currentFile.IsRename = true
	currentFile.NewName = strings.Fields(line)[2]
}

func parseOldMode(currentFile *FileDiff, line string) {
	currentFile.OldMode = strings.Fields(line)[2]
}

func parseNewMode(currentFile *FileDiff, line string) {
	currentFile.NewMode = strings.Fields(line)[2]
}
