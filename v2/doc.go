package testmtx

// Doc is a simple document interface.
type Doc interface {
	GetSheetNames() ([]string, error)
	GetSheet(sheetName string) (DocSheet, error)
}

// DocSheet is a simple sheet interface.
type DocSheet interface {
	Name() string
	Rows() []DocRow
	Value(row, clm int) string
}

// DocRow is a simple row interface.
type DocRow interface {
	Number() int
	Value(clm int) string
}
