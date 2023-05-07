package endpoint

import (
	"context"
	"github.com/zander-84/seagull/transport"
)

type Endpoint func(ctx context.Context, request interface{}, codec Codec, errorEncoder ErrorEncoder, recoverEncoder RecoverEncoder) (response interface{}, err error)

type HandlerFunc func(ctx context.Context, request interface{}) (response interface{}, err error)

type ProxyEndpoint func(p transport.Kind, fullPath Path, e HandlerFunc)

type ErrorEncoder func(ctx context.Context, p transport.Kind, err error) error

func WrapError(e map[transport.Kind]func(ctx context.Context, err error) error) ErrorEncoder {
	if e == nil {
		return nil
	}
	return func(ctx context.Context, p transport.Kind, err error) error {
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

func WrapRecover(rec map[transport.Kind]func(ctx context.Context, err *error)) func(ctx context.Context, p transport.Kind) RecoverEncoder {
	if rec == nil {
		return nil
	}

	return func(ctx context.Context, p transport.Kind) RecoverEncoder {
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

type Codecs map[transport.Kind]Codec

type Codec struct {
	Dec Dec
	Enc Enc
}

type Dec func(ctx context.Context, in any) (out any, err error)
type Enc func(ctx context.Context, in any) (out any, err error)
