package grpcerror

import (
	pb "main/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ExternalGroup = "external"
	InternalGroup = "internal"

	MessageDescriptionSizeLimit = 2000
)

type Error struct {
	Group        string
	Code         string
	Description  string
	InternalCode codes.Code
}

func NewErorFromProto(err error) *Error {
	st := status.Convert(err)

	if st.Code() != codes.Unavailable && st.Code() != codes.Internal {
		return nil
	}

	details := st.Details()
	if len(details) == 0 {
		return nil
	}

	switch detail := details[0].(type) { // ожидаем только один элемент в массиве
	case *pb.Error:
		return &Error{
			Group:        detail.GetGroup(),
			Code:         detail.GetCode(),
			Description:  detail.GetDescription(),
			InternalCode: codes.Code(detail.GetInternalCode()),
		}
	default:
		return nil
	}
}

func NewExternal() *Error {
	return &Error{
		Group:        ExternalGroup,
		InternalCode: codes.Unavailable,
	}
}

func NewInternal() *Error {
	return &Error{
		Group:        InternalGroup,
		InternalCode: codes.Internal,
	}
}

func (e *Error) WithDescription(desc string) *Error {
	e.Description = desc

	return e
}

func (e *Error) WithCode(code string) *Error {
	e.Code = code

	return e
}

func (e *Error) Error() string {
	return e.Code
}

func (e *Error) GRPCStatus() *status.Status {
	details := pb.Error{
		Group:        e.Group,
		Code:         e.Code,
		Description:  e.Description,
		InternalCode: uint32(e.InternalCode),
	}

	if e.size() >= MessageDescriptionSizeLimit {
		details.Description = "ERROR MESSAGE SIZE LIMIT EXCEEDED"
	}

	status, _ := status.New(e.InternalCode, "").WithDetails(&details)

	return status
}

func (e *Error) size() int {
	return (len(e.Code) + len(e.Description) + len(e.Group) + 4)
}
