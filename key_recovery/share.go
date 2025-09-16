package keyrecovery

// A Share represents data given to a guardian.
//
// All fields are hex strings without the leading 0x
type Share struct {
	Point        string // point of evaluation
	SpendingEval string // spending polynomial evaluated on Point
	ViewingEval  string // viewing polynomial evaluated on Point
}
