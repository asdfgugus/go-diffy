package godiffy

type Diff struct {
	Files []*FileDiff
}

type FileStatus int

type FileDiff struct {
	Header          string
	OldHash         string
	NewHash         string
	SimilarityIndex string
	OldPath         string
	NewPath         string
	OldName         string
	NewName         string
	OldMode         string
	NewMode         string
	Status          FileStatus
	Hunks           []*Hunk
}

type Hunk struct {
	OldStart     int
	NewStart     int
	OldLineCount int
	NewLineCount int
	Lines        []*HunkLine
}

type HunkLineKind int

type HunkLine struct {
	Type    HunkLineKind
	Content string
}
