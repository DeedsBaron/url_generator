package errz

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Error(t *testing.T) {
	t.Run("base use case", func(t *testing.T) {
		err := Error(Unavailable, "msg1")
		require.Equal(t, err.Error(), "code = UNAVAILABLE desc = msg1")
		err = WrapC(err, InvalidArgument, "msg2")
		require.Equal(t, err.Error(), "code = INVALID_ARGUMENT desc = msg2: msg1")
		err = errors.Unwrap(err)
		require.Equal(t, err.Error(), "code = UNAVAILABLE desc = msg1")
		err = errors.Unwrap(err)
		require.Nil(t, err)
	})

	t.Run("wrap std err", func(t *testing.T) {
		err := errors.New("std err")
		err = Wrap(err, "wrapped")
		require.Equal(t, err.Error(), "code = INTERNAL desc = wrapped: std err")
		err = errors.Unwrap(err)
		require.Equal(t, err.Error(), "std err")
		err = errors.Unwrap(err)
		require.Nil(t, err)
	})

	t.Run("status code", func(t *testing.T) {
		err := Error(Internal, "err")
		var e StatusError
		require.ErrorAs(t, err, &e)
		require.Equal(t, codes.Internal, e.GRPCStatus().Code())
		require.Equal(t, Internal, e.Code())

		err = Error(Unimplemented, "err")
		require.ErrorAs(t, err, &e)
		require.Equal(t, codes.Unimplemented, e.GRPCStatus().Code())
		require.Equal(t, Unimplemented, e.Code())
	})

	t.Run("proto message", func(t *testing.T) {
		msg := "something went wrong"
		err := Error(Internal, msg)
		protoStatus, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, msg, protoStatus.Message())
	})

	t.Run("std errors.Is", func(t *testing.T) {
		var ErrNotFound = Error(NotFound, "not found")
		err := Wrap(ErrNotFound, "package")
		err2 := Wrap(err, "another package")

		require.True(t, errors.Is(err2, ErrNotFound))
	})

	t.Run("custom error code", func(t *testing.T) {
		DBError := NewCode("DB_ERROR", Internal)
		ErrDB := Errorf(DBError, "db unavailable")
		err := Wrap(ErrDB, "package")

		require.Equal(t, err.Error(), "code = DB_ERROR desc = package: db unavailable")
		var e StatusError
		require.ErrorAs(t, err, &e)
		require.Equal(t, DBError, e.Code())
		require.Equal(t, codes.Internal, e.GRPCStatus().Code())
		require.Equal(t, "package: db unavailable", e.GRPCStatus().Message())
		require.Error(t, e.GRPCStatus().Err())
	})

	t.Run("unwrap code", func(t *testing.T) {
		err := Error(NotFound, "not found")
		wrapped := WrapC(err, InvalidArgument, "wrap error")
		require.Equal(t, InvalidArgument, UnwrapCode(wrapped))
	})
	t.Run("wrapf and wrapfc with non errz error", func(t *testing.T) {
		err := errors.New("some error")
		wrapped := Wrapf(err, "a %d", 1)
		require.Equal(t, "code = INTERNAL desc = a 1: some error", wrapped.Error())
		wrapped = WrapfC(err, NotFound, "a %d", 2)
		require.Equal(t, "code = NOT_FOUND desc = a 2: some error", wrapped.Error())
	})

	t.Run("error message without code", func(t *testing.T) {
		msg := "something went wrong"
		err := Error(Internal, msg)
		require.Equal(t, msg, ErrorWithoutCode(err))
	})

	t.Run("wrap and wrapf determines the code from well-known errors", func(t *testing.T) {
		err := Wrap(errors.New("oops"), "something happened")
		var e StatusError
		require.ErrorAs(t, err, &e)
		require.Equal(t, Internal, e.Code())

		err = Wrapf(context.Canceled, "something happened with id=%d", 2)
		require.ErrorAs(t, err, &e)
		require.Equal(t, Canceled, e.Code())

		err = Wrap(context.DeadlineExceeded, "lvl 1")
		require.ErrorAs(t, err, &e)
		require.Equal(t, DeadlineExceeded, e.Code())

		err = Wrap(WrapC(err, FailedPrecondition, "lvl 2"), "lvl 3")
		require.ErrorAs(t, err, &e)
		require.Equal(t, FailedPrecondition, e.Code())

		err = Wrap(fmt.Errorf("wrap: %w", context.Canceled), "wrap")
		require.ErrorAs(t, err, &e)
		require.Equal(t, Canceled, e.Code())
	})
}

func Test_StatusErrorWrap(t *testing.T) {
	t.Run("wrap status error", func(t *testing.T) {
		err := status.Error(codes.Canceled, "canceled error")
		err = Wrap(err, "wrapped canceled err")
		require.Equal(t, Canceled, UnwrapCode(err))

		err = status.Error(codes.InvalidArgument, "invalid argument error")
		err = Wrap(err, "wrapped invalid argument err")
		require.Equal(t, InvalidArgument, UnwrapCode(err))
	})

	t.Run("wrap unknown code error", func(t *testing.T) {
		err := status.Error(codes.OK, "ok error (WAT)???")
		err = Wrap(err, "wrapped canceled err")
		require.Equal(t, Internal, UnwrapCode(err))
	})
}

func Test_errz_GRPCStatus(t *testing.T) {
	for codeStr, _ := range strToCode {
		t.Run(codeStr, func(t *testing.T) {
			err := Error(code(codeStr), "error happens")
			var e StatusError
			errors.As(err, &e)

			actualStatus, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, e.GRPCStatus(), actualStatus)
			// there must be an error, otherwise there will be panic in the scratch
			require.Error(t, actualStatus.Err())
		})
	}
}

func Test_errz_Builder(t *testing.T) {
	err := New().Messagef("a 1").Code(NotFound).Params(map[string]any{
		"posting": 1,
	}).Build()
	require.Equal(t, "code = NOT_FOUND desc = a 1", err.Error())
	err = New().Wrap(err, "a 2").Build()
	require.Equal(t, "code = NOT_FOUND desc = a 2: a 1", err.Error())

	err = New().Wrap(context.Canceled, "b 3").Build()
	require.Equal(t, "code = CANCELED desc = b 3: context canceled", err.Error())
}
