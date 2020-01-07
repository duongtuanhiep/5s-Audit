package enum

//ItemAnswer enum
type ItemAnswer int

// itemAnswer : [SORT,ORDER,SHINE,STANDARDIZE,SUSTAIN]
const (
	HORRIBLE ItemAnswer = iota
	BAD
	AVERAGE
	GOOD
	EXCELLENT
)
