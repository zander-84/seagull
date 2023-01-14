package endpoint

import (
	"context"
)

type Endpoint func(ctx context.Context, request interface{}, codec Codec, errorEncoder ErrorEncoder, recoverEncoder RecoverEncoder) (response interface{}, err error)

type HandlerFunc func(ctx context.Context, request interface{}) (response interface{}, err error)

type ProxyEndpoint func(p Kind, fullPath Path, e HandlerFunc)

type ErrorEncoder func(ctx context.Context, p Kind, err error) error

func WrapError(e map[Kind]func(ctx context.Context, err error) error) ErrorEncoder {
	if e == nil {
		return nil
	}
	return func(ctx context.Context, p Kind, err error) error {
		for k, v := range e {
			if k == p {
				if v == nil {
					break
				}
				return v(ctx, err)
			}
		}
		return err
	}
}

type RecoverEncoder func(ctx context.Context, err *error)

func WrapRecover(rec map[Kind]func(ctx context.Context, err *error)) func(ctx context.Context, p Kind) RecoverEncoder {
	if rec == nil {
		return nil
	}

	return func(ctx context.Context, p Kind) RecoverEncoder {
		for k, v := range rec {
			if k == p {
				if v == nil {
					break
				}
				return v
			}
		}

		return nil
	}
}

type Codecs map[Kind]Codec

type Codec struct {
	Dec Dec
	Enc Enc
}

type Dec func(ctx context.Context, in any) (out any, err error)
type Enc func(ctx context.Context, in any) (out any, err error)
