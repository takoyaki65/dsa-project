package requeststatus

type RequestResultState int64

// refer to docs/README.md
const (
	AC RequestResultState = iota
	WA
	TLE
	MLE
	RE
	CE
	OLE
	IE
	FN
	Judging
	WJ
)
