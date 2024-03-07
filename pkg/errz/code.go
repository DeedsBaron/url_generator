package errz

import "google.golang.org/grpc/codes"

// Code for errz
type Code interface {
	GRPCCode() codes.Code
	String() string
}

type code string

func (c code) String() string {
	return string(c)
}

func (c code) GRPCCode() codes.Code {
	return strToCode[string(c)]
}

var strToCode = map[string]codes.Code{
	"CANCELED":            codes.Canceled,
	"UNKNOWN":             codes.Unknown,
	"INVALID_ARGUMENT":    codes.InvalidArgument,
	"DEADLINE_EXCEEDED":   codes.DeadlineExceeded,
	"NOT_FOUND":           codes.NotFound,
	"ALREADY_EXISTS":      codes.AlreadyExists,
	"PERMISSION_DENIED":   codes.PermissionDenied,
	"RESOURCE_EXHAUSTED":  codes.ResourceExhausted,
	"FAILED_PRECONDITION": codes.FailedPrecondition,
	"ABORTED":             codes.Aborted,
	"OUT_OF_RANGE":        codes.OutOfRange,
	"UNIMPLEMENTED":       codes.Unimplemented,
	"INTERNAL":            codes.Internal,
	"UNAVAILABLE":         codes.Unavailable,
	"DATA_LOSS":           codes.DataLoss,
	"UNAUTHENTICATED":     codes.Unauthenticated,
}

var statusCodeToCode = map[codes.Code]Code{
	codes.Canceled:           Canceled,
	codes.Unknown:            Unknown,
	codes.InvalidArgument:    InvalidArgument,
	codes.DeadlineExceeded:   DeadlineExceeded,
	codes.NotFound:           NotFound,
	codes.AlreadyExists:      AlreadyExists,
	codes.PermissionDenied:   PermissionDenied,
	codes.ResourceExhausted:  ResourceExhausted,
	codes.FailedPrecondition: FailedPrecondition,
	codes.Aborted:            Aborted,
	codes.OutOfRange:         OutOfRange,
	codes.Unimplemented:      Unimplemented,
	codes.Internal:           Internal,
	codes.Unavailable:        Unavailable,
	codes.DataLoss:           DataLoss,
	codes.Unauthenticated:    Unauthenticated,
}

// errz error codes.
const (
	Canceled           code = "CANCELED"
	Unknown            code = "UNKNOWN"
	InvalidArgument    code = "INVALID_ARGUMENT"
	DeadlineExceeded   code = "DEADLINE_EXCEEDED"
	NotFound           code = "NOT_FOUND"
	AlreadyExists      code = "ALREADY_EXISTS"
	PermissionDenied   code = "PERMISSION_DENIED"
	ResourceExhausted  code = "RESOURCE_EXHAUSTED"
	FailedPrecondition code = "FAILED_PRECONDITION"
	Aborted            code = "ABORTED"
	OutOfRange         code = "OUT_OF_RANGE"
	Unimplemented      code = "UNIMPLEMENTED"
	Internal           code = "INTERNAL"
	Unavailable        code = "UNAVAILABLE"
	DataLoss           code = "DATA_LOSS"
	Unauthenticated    code = "UNAUTHENTICATED"
)

// GetCode return Code
func GetCode(cd codes.Code) Code {
	if c, ok := statusCodeToCode[cd]; ok {
		return c
	}
	return Internal
}

type customCode struct {
	name   string
	analog code
}

func (c customCode) GRPCCode() codes.Code {
	return c.analog.GRPCCode()
}

func (c customCode) String() string {
	return c.name
}

// NewCode creates a new error Code.
func NewCode(name string, analog code) Code {
	return &customCode{
		name:   name,
		analog: analog,
	}
}
