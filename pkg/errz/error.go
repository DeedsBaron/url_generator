package errz

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/status"
)

// StatusError represents an error with a status code.
type StatusError interface {
	error
	Code() Code
	GRPCStatus() *status.Status
	Message() string
}

type errz struct {
	params map[string]any
	code   Code
	msg    string
	cause  error
}

// GRPCStatus returns the gRPC status.
func (e *errz) GRPCStatus() *status.Status {
	return status.New(e.code.GRPCCode(), e.msg)
}

func (e *errz) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	for k, v := range e.params {
		switch val := v.(type) {
		case int:
			encoder.AddInt(k, val)
		case bool:
			encoder.AddBool(k, val)
		case string:
			encoder.AddString(k, val)
		case int64:
			encoder.AddInt64(k, val)
		case zapcore.ObjectMarshaler:
			err := encoder.AddObject(k, val)
			if err != nil {
				return err
			}
		default:
			err := encoder.AddReflected(k, val)
			if err != nil {
				return err
			}
		}
	}
	encoder.AddString("message", e.Message())
	return nil
}

// Code returns an internal error Code.
func (e *errz) Code() Code {
	return e.code
}

// Error implements the built-in error interface.
func (e *errz) Error() string {
	return fmt.Sprintf("code = %s desc = %s", e.code, e.msg)
}

// UnwrapCode determines the error Code.
func UnwrapCode(err error) Code {
	if e, ok := err.(interface {
		Code() Code
	}); ok {
		return e.Code()
	}

	if e, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
		return GetCode(e.GRPCStatus().Code())
	}

	if errors.Is(err, context.Canceled) {
		return Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return DeadlineExceeded
	}

	return Internal
}

// Error returns new errz error
func Error(code Code, msg string) error {
	return &errz{
		code: code,
		msg:  msg,
	}
}

// Errorf returns new errz error with formatting
func Errorf(code Code, format string, args ...interface{}) error {
	return &errz{
		code: code,
		msg:  fmt.Sprintf(format, args...),
	}
}

// Wrap wraps err with message
func Wrap(err error, msg string) error {
	return &errz{
		code:  UnwrapCode(err),
		msg:   fmt.Sprintf("%s: %s", msg, ErrorWithoutCode(err)),
		cause: err,
	}
}

// Wrapf wraps err with format
func Wrapf(err error, format string, args ...interface{}) error {
	return &errz{
		code:  UnwrapCode(err),
		msg:   fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), ErrorWithoutCode(err)),
		cause: err,
	}
}

// WrapC wraps err with message and new code
func WrapC(err error, code Code, msg string) error {
	return &errz{
		code:  code,
		msg:   fmt.Sprintf("%s: %s", msg, ErrorWithoutCode(err)),
		cause: err,
	}
}

// WrapfC wraps err with format and new code
func WrapfC(err error, code Code, format string, args ...interface{}) error {
	return &errz{
		code:  code,
		msg:   fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), ErrorWithoutCode(err)),
		cause: err,
	}
}

// Unwrap if err is errz get wrapped err. May return nil
func (e *errz) Unwrap() error {
	return e.cause
}

// Cause if err is errz get wrapped err. May return nil
func (e *errz) Cause() error {
	return e.cause
}

// Builder errz builder
type Builder struct {
	params map[string]any
	code   Code
	msg    string
	cause  error
}

// New creates new builder
func New() *Builder {
	return &Builder{
		code: Internal,
	}
}

// Message set builder message
func (b *Builder) Message(msg string) *Builder {
	b.msg = msg
	return b
}

// Messagef set builder message with format
func (b *Builder) Messagef(msg string, args ...any) *Builder {
	b.msg = fmt.Sprintf(msg, args...)
	return b
}

// Wrap wraps the error.
func (b *Builder) Wrap(err error, msg string) *Builder {
	b.code = UnwrapCode(err)
	b.msg = fmt.Sprintf("%s: %v", msg, ErrorWithoutCode(err))
	b.cause = err

	return b
}

// Code set code
func (b *Builder) Code(c Code) *Builder {
	b.code = c
	return b
}

// Params set params
func (b *Builder) Params(p map[string]any) *Builder {
	b.params = p
	return b
}

// Build build error
func (b *Builder) Build() error {
	return &errz{
		params: b.params,
		code:   b.code,
		msg:    b.msg,
		cause:  b.cause,
	}
}

// Message get message from error
func (e *errz) Message() string {
	return e.msg
}

// ErrorWithoutCode get message from error without code
func ErrorWithoutCode(err error) string {
	if err == nil {
		return "nil error"
	}
	if e, ok := err.(interface {
		Message() string
	}); ok {
		return e.Message()
	}
	return err.Error()
}
