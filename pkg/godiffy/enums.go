package godiffy

const (
	FileStatusAdded FileStatus = iota
	FileStatusDeleted
	FileStatusModified
	FileStatusRenamed
	FileStatusCopied
	FileStatusUnknown
)

const (
	LineAdded HunkLineKind = iota
	LineDeleted
	LineContext
)
