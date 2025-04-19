package godiffy

type Diff struct {
	Files []*FileDiff
}

type FileStatus int

type FileDiff struct {
	Header          string
	IsNew           bool
	IsDeleted       bool
	IsRename        bool
	IsCopy          bool
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
}

type Hunk struct {
	OldStart int
	NewStart int
	OldLines int
	NewLines int
	Lines    []HunkLine
}

type HunkLineKind int

type HunkLine struct {
	Type    HunkLineKind
	Content string
}
