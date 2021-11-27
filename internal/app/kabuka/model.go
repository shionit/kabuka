package kabuka

import "golang.org/x/xerrors"

// Option execution parameters
type Option struct {
	Symbol string
	//ShowDetail        bool
	Format OutputFormatType // console(default) or json
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
)

func ParseOutputFormat(s string) (OutputFormatType, error) {
	switch s {
	case "", "text":
		return OutputFormatTypeText, nil
	case "json":
		return OutputFormatTypeJson, nil
	}
	return "", xerrors.Errorf("Unsupported format: %s", s)
}
