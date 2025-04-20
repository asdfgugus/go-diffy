package godiffy

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(input string) (*Diff, error) {
	resultDiff := &Diff{}
	var currentFile *FileDiff
	var currentHunk *Hunk
	isHeader := true
	isHunk := false

	for line := range strings.Lines(input) {
		if currentFile == nil && !strings.HasPrefix(line, "diff --git") {
			continue
		}
		if strings.HasPrefix(line, "di") { // diff --git a/foo.txt b/foo.txt
			currentFile = parseNewFileDiff(resultDiff, line)
			isHeader = true
			isHunk = false
			continue
		}
		if strings.HasPrefix(line, "@") {
			var err error
			currentHunk, err = parseHunk(currentFile, line)
			if err != nil {
				return nil, err
			}
			isHeader = false
			isHunk = true
			continue
		}

		if isHeader {
			switch {
			case strings.HasPrefix(line, "i"): // index abc123..def456 100644
				err := parseMetadata(currentFile, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "--- "): // --- a/foo.txt
				err := parseOldFilenameMarker(currentFile, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "+++ "): // +++ b/foo.txt
				err := parseNewFilenameMarker(currentFile, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "new f"): // new file mode 100644
				err := parseNewFileMode(currentFile, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "de"): // deleted file mode 100644
				err := parseDeletedFileMode(currentFile, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "r"): // rename from old_name.txt / rename to new_name.txt
				if strings.HasPrefix(line, "rename f") {
					err := parseRenameFrom(currentFile, line)
					if err != nil {
						return nil, err
					}
				} else {
					err := parseRenameTo(currentFile, line)
					if err != nil {
						return nil, err
					}
				}
			case strings.HasPrefix(line, "o"): // old mode 100755
				err := parseOldMode(currentFile, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "new m"): // new mode 100644
				err := parseNewMode(currentFile, line)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("failed to parse line: %s", line)
			}
		}

		if isHunk {
			switch {
			case strings.HasPrefix(line, "+"): // +added line
				err := parseAddLine(currentHunk, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "-"): // -removed line
				err := parseDeleteLine(currentHunk, line)
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, " "): //  context line
				err := parseContextLine(currentHunk, line)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("failed to parse line: %s", line)
			}
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

func parseHunk(currentFile *FileDiff, line string) (*Hunk, error) {
	var err error
	hunk := &Hunk{}
	parts := strings.Fields(line) // ["@@","-1,3","+1,4","@@"]
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid hunk format: %s", line)
	}
	oldParts := strings.Split(parts[1], ",") // ["-1","3"]
	newParts := strings.Split(parts[2], ",")
	if len(oldParts) != 2 && len(newParts) != 2 {
		return nil, fmt.Errorf("invalid hunk format: %s", line)
	}
	hunk.OldStart, err = strconv.Atoi(strings.TrimPrefix(oldParts[0], "-")) // -1
	if err != nil {
		return nil, fmt.Errorf("failed to parse old start line %s: %w", line, err)
	}
	hunk.NewStart, err = strconv.Atoi(strings.TrimPrefix(newParts[0], "+")) // +1
	if err != nil {
		return nil, fmt.Errorf("failed to parse new start line %s: %w", line, err)
	}
	hunk.OldLineCount, err = strconv.Atoi(oldParts[1]) // 3
	if err != nil {
		return nil, fmt.Errorf("failed to parse old lines %s: %w", line, err)
	}
	hunk.NewLineCount, err = strconv.Atoi(newParts[1]) // 4
	if err != nil {
		return nil, fmt.Errorf("failed to parse new lines %s: %w", line, err)
	}

	currentFile.Hunks = append(currentFile.Hunks, hunk)
	return hunk, nil
}

func parseOldFilenameMarker(currentFile *FileDiff, line string) error {
	parts := strings.SplitN(line, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid filename format: %s", line)
	}
	currentFile.OldPath = strings.TrimSpace(parts[1])
	return nil
}

func parseNewFilenameMarker(currentFile *FileDiff, line string) error {
	parts := strings.SplitN(line, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid filename format: %s", line)
	}
	currentFile.NewPath = strings.TrimSpace(parts[1])
	return nil
}

func parseAddLine(hunk *Hunk, line string) error {
	if hunk == nil {
		return fmt.Errorf("failed to parse add line: hunk is nil")
	}
	hunkLine := HunkLine{
		Type:    HunkLineAdded,
		Content: strings.TrimPrefix(line, "+"),
	}
	hunk.Lines = append(hunk.Lines, hunkLine)
	return nil
}

func parseDeleteLine(hunk *Hunk, line string) error {
	if hunk == nil {
		return fmt.Errorf("failed to parse delete line: hunk is nil")
	}
	hunkLine := HunkLine{
		Type:    HunkLineDeleted,
		Content: strings.TrimPrefix(line, "-"),
	}
	hunk.Lines = append(hunk.Lines, hunkLine)
	return nil
}

func parseContextLine(hunk *Hunk, line string) error {
	if hunk == nil {
		return fmt.Errorf("failed to parse context line: hunk is nil")
	}
	hunkLine := HunkLine{
		Type:    HunkLineContext,
		Content: strings.TrimPrefix(line, " "),
	}
	hunk.Lines = append(hunk.Lines, hunkLine)
	return nil
}

func parseMetadata(currentFile *FileDiff, line string) error {
	parts := strings.Fields(line) // ["index","abc123..def456","100644"]
	if len(parts) < 2 {
		return fmt.Errorf("invalid metadata format: %s", line)
	}
	hashes := strings.SplitN(parts[1], "..", 2)
	if len(hashes) != 2 {
		return fmt.Errorf("invalid hash format: %s", line)
	}
	currentFile.OldHash, currentFile.NewHash = strings.TrimSpace(hashes[0]), strings.TrimSpace(hashes[1])
	if len(parts) > 2 {
		currentFile.NewMode = strings.TrimSpace(parts[2])
	}
	return nil
}

func parseNewFileMode(currentFile *FileDiff, line string) error {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return fmt.Errorf("invalid new file mode format: %s", line)
	}
	currentFile.NewMode = strings.TrimSpace(parts[3])
	currentFile.Status = FileStatusNew
	return nil
}

func parseDeletedFileMode(currentFile *FileDiff, line string) error {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return fmt.Errorf("invalid deleted file mode format: %s", line)
	}
	currentFile.OldMode = strings.TrimSpace(parts[3])
	currentFile.Status = FileStatusDeleted
	return nil
}

func parseRenameFrom(currentFile *FileDiff, line string) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("invalid rename from format: %s", line)
	}
	currentFile.OldName = strings.TrimSpace(parts[2])
	currentFile.Status = FileStatusRenamed
	return nil
}

func parseRenameTo(currentFile *FileDiff, line string) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("invalid rename to format: %s", line)
	}
	currentFile.NewName = strings.TrimSpace(parts[2])
	currentFile.Status = FileStatusRenamed
	return nil
}

func parseOldMode(currentFile *FileDiff, line string) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("invalid old mode format: %s", line)
	}
	currentFile.OldMode = strings.TrimSpace(parts[2])
	return nil
}

func parseNewMode(currentFile *FileDiff, line string) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("invalid new mode format: %s", line)
	}
	currentFile.NewMode = strings.TrimSpace(parts[2])
	return nil
}
