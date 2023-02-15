package think

import (
	"fmt"
	"github.com/zander-84/seagull/contrib/codec"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type OptsResp func(e *Response)

func optsResp(err *Response, opts ...OptsResp) {
	for _, v := range opts {
		v(err)
	}
}

func SetBizCode(biz string) OptsResp {
	return func(e *Response) {
		e.BizCode = biz
	}
}

func SetMessage(message string) OptsResp {
	return func(e *Response) {
		e.Message = message
	}
}

func SetMetadata(metadata map[string]string) OptsResp {
	return func(e *Response) {
		e.Metadata = metadata
	}
}

type Response struct {
	Code     Code
	BizCode  string
	Message  string
	Metadata map[string]string // 比如 page info
	Data     interface{}
}

func NewResponse(code Code, bizCode string, message string, metadata map[string]string, data interface{}) *Response {
	return &Response{
		Code:     code,
		BizCode:  bizCode,
		Message:  message,
		Metadata: metadata,
		Data:     data,
	}
}

func (r *Response) ToJson() string {
	res, _ := codec.GetCodec(codec.Json).Marshal(r)
	return string(res)
}

func (r *Response) ToProtoMessage(body proto.Message) proto.Message {
	responseGrpc := &ResponseGrpc{
		Code:     r.Code,
		BizCode:  r.BizCode,
		Message:  r.Message,
		Metadata: r.Metadata,
		Data:     nil,
	}
	if body != nil {
		anyErr, _ := anypb.New(body)
		responseGrpc.Data = anyErr
	}

	return responseGrpc
}
func (r *Response) Unmarshal(in []byte, body any, name string) error {
	if name == codec.Proto {
		body2, ok := body.(proto.Message)
		if !ok {
			return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", body)
		}
		return r.UnmarshalPorto(in, body2)
	} else if name == codec.Json {
		return r.UnmarshalJson(in, body)
	} else {
		return fmt.Errorf("error codec name :%s", name)
	}
}

func (r *Response) UnmarshalJson(in []byte, body any) error {
	if err := codec.GetCodec(codec.Json).Unmarshal(in, r); err != nil {
		return err
	}
	if r.Code < Code_Min {
		return fmt.Errorf("code less than min code: %d", r.Code)
	}
	body = r.Data
	return nil
}

func (r *Response) UnmarshalPorto(in []byte, body proto.Message) error {
	pb := new(ResponseGrpc)
	if err := codec.GetCodec(codec.Proto).Unmarshal(in, pb); err != nil {
		return err
	}

	if pb.Code < Code_Min {
		return fmt.Errorf("code less than min code: %d", r.Code)
	}

	if pb.Data != nil && body != nil {
		if err := pb.Data.UnmarshalTo(body); err != nil {
			return err
		}
	}

	r.Code = pb.Code
	r.BizCode = pb.BizCode
	r.Message = pb.Message
	r.Data = body

	return nil
}

//----

func NewSuccessResp(data any, opts ...OptsResp) *Response {
	out := NewResponse(Code_Success, "", "", nil, data)
	optsResp(out, opts...)
	return out
}
