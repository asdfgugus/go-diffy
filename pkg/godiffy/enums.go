package godiffy

const (
	FileStatusNew FileStatus = iota
	FileStatusDeleted
	FileStatusModified
	FileStatusRenamed
	FileStatusCopied
	FileStatusUnknown
)

const (
	HunkLineAdded HunkLineKind = iota
	HunkLineDeleted
	HunkLineContext
)
