package kabuka

import "golang.org/x/xerrors"

// Option execution parameters
type Option struct {
	Symbol string
	//ShowDetail        bool
	Format OutputFormatType // text(default) or json or csv
	//OutputColumnsType string // default, all, individual
	//OutputColumns     []string
}

type Kabuka struct {
	Option
}

type OutputFormatType string

const (
	OutputFormatTypeText OutputFormatType = "text"
	OutputFormatTypeJson OutputFormatType = "json"
	OutputFormatTypeCsv  OutputFormatType = "csv"
)

func ParseOutputFormat(s string) (OutputFormatType, error) {
	switch s {
	case "", "text":
		return OutputFormatTypeText, nil
	case "json":
		return OutputFormatTypeJson, nil
	case "csv":
		return OutputFormatTypeCsv, nil
	}
	return "", xerrors.Errorf("Unsupported format: %s", s)
}
