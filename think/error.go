package think

import (
	"errors"
	"fmt"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type Error struct {
	*Response
	cause error
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Errors() string {
	return fmt.Sprintf("error: code = %d bizcode = %s data = %v message = %s metadata = %v cause = %v", e.Code, e.BizCode, e.Data, e.Message, e.Metadata, e.cause)
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Message == e.Message
	}
	return false
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := CloneErr(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := CloneErr(e)
	err.Metadata = md
	return err
}

// NewErr returns an error object for the code, message.
func NewErr(code Code, bizCode string, message string) *Error {
	return &Error{
		Response: &Response{
			Code:    code,
			BizCode: bizCode,
			Message: message,
		},
	}
}

// CloneErr deep clone error to a new error.
func CloneErr(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		cause: err.cause,
		Response: &Response{
			Code:     err.Code,
			Data:     err.Data,
			Message:  err.Message,
			Metadata: metadata,
		},
	}
}

func (e *Error) ToGrpcErr() error {
	// return status.Error(codes.Code(e.Code), e.Response.ToJson())

	message := e.ToProtoMessage(nil)
	anyBody, _ := anypb.New(message)

	return status.ErrorProto(&spb.Status{
		Code:    int32(e.Code),
		Message: "",
		Details: []*anypb.Any{
			anyBody,
		},
	})
}

func GrpcErr2ThinkErr(err error) *Error {
	gs, ok := status.FromError(err)
	if !ok {
		return NewErr(Code_Undefined, "", err.Error())
	}

	if int32(gs.Code()) < int32(Code_Min) {
		return NewErr(Code_SystemSpaceError, "", gs.Message())
	}

	_proto := gs.Proto()
	if _proto == nil || len(_proto.Details) < 1 {
		return NewErr(Code_SystemSpaceError, "", "rpc lost detail")
	}

	resp := new(Response)
	if err := resp.UnmarshalPorto(_proto.Details[0].Value, nil); err != nil {
		return NewErr(Code_SystemSpaceError, "", err.Error())
	}

	return &Error{
		Response: resp,
	}
}

// GetCode returns the http code for an error.
// It supports wrapped errors.
func GetCode(err error) Code {
	if err == nil {
		return Code_Success
	}
	return FromError(err).Code
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	return NewErr(Code_Undefined, "", err.Error())
}

//-----------

var RecordNotFound = NewErr(Code_NotFound, "", "record not found")
var RecordExist = NewErr(Code_Repeat, "", "record already exist")
var SystemError = NewErr(Code_SystemSpaceError, "", "system error")
var UnImpl = NewErr(Code_UnImpl, "", "UnImpl")

//-----------

func IsErrSystemSpace(err error) bool {
	return GetCode(err) == Code_SystemSpaceError
}

func IsErrUndefined(err error) bool {
	return GetCode(err) == Code_Undefined
}
func IsErrNotFound(err error) bool {
	return GetCode(err) == Code_NotFound
}

//-----------

func ErrUndefined(message string, opts ...OptsResp) *Error {
	err := NewErr(Code_Undefined, "", message)
	optsResp(err.Response, opts...)
	return err
}

func ErrSystemSpace(message string, opts ...OptsResp) *Error {
	err := NewErr(Code_SystemSpaceError, "", message)
	optsResp(err.Response, opts...)
	return err
}
func ErrTooManyRequests(message string, opts ...OptsResp) *Error {
	err := NewErr(Code_TooManyRequests, "", message)
	optsResp(err.Response, opts...)
	return err
}
func ErrRecordNotFound(message string, opts ...OptsResp) *Error {
	err := NewErr(Code_NotFound, "", message)
	optsResp(err.Response, opts...)
	return err
}

func ErrAlert(message string, opts ...OptsResp) *Error {
	err := NewErr(Code_AlterError, "", message)
	optsResp(err.Response, opts...)
	return err
}

func ErrType(message string, opts ...OptsResp) *Error {
	err := NewErr(Code_TypeError, "", message)
	optsResp(err.Response, opts...)
	return err
}
