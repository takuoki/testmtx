package testmtx

// TODO: publicである必要があるか

type Doc interface {
	GetSheetNames() ([]string, error)
	GetSheet(sheetName string) (DocSheet, error)
}

// DocSheet is a simple sheet interface.
type DocSheet interface {
	Rows() []DocRow
	Value(row, clm int) string
}

// DocRow is a simple row interface.
type DocRow interface {
	Value(clm int) string
}
