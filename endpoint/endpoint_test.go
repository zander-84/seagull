package endpoint

import (
	"context"
	"fmt"
	"testing"
)

func TestChain1(t *testing.T) {
	e := Chain(
		annotate("first"),
	)(myEndpoint)

	if _, err := e(context.Background(), struct{}{}); err != nil {
		panic(err)
	}
	/*
		annotate("first"):
			 func(next Endpoint) Endpoint {
				return func(ctx context.Context, request interface{}) (interface{}, error) {
					fmt.Println("first", "pre")
					return next(ctx, request)
					fmt.Println("first", "post")
				}
			}

		Chain(annotate("first")):
			 func(next Endpoint) Endpoint {
				return func(ctx context.Context, request interface{}) (interface{}, error) {
					fmt.Println("first", "pre")
					return next(ctx, request)
					fmt.Println("first", "post")
				}
			}

		Chain(annotate("first"))(myEndpoint):
			 func(ctx context.Context, request interface{}) (interface{}, error) {
					fmt.Println("first", "pre")
					return myEndpoint(ctx, request)
					fmt.Println("first", "post")
			}

		Chain(annotate("first"))(myEndpoint)(context.Background(), struct{}{}):
			func(ctx context.Context, request interface{}) (interface{}, error) {
					fmt.Println("first", "pre")
					return myEndpoint(ctx, request)
					fmt.Println("first", "post")
			}(context.Background(), struct{}{})
	*/
}

func TestChain(t *testing.T) {
	e := Chain(
		annotate("first"),
		annotate("second"),
	)(myEndpoint)

	if _, err := e(context.Background(), struct{}{}); err != nil {
		panic(err)
	}
	/*
		1.参数解析
		annotate("first"):
			next:= func(next Endpoint) Endpoint {
				return func(ctx context.Context, request interface{}) (interface{}, error) {
					fmt.Println("first", "pre")
					return next(ctx, request)
					fmt.Println("first", "post")
				}
			}
		annotate("second"):
			next:= func(next Endpoint) Endpoint {
				return func(ctx context.Context, request interface{}) (interface{}, error) {
					fmt.Println("second", "pre")
					return next(ctx, request)
					fmt.Println("second", "post")
				}
			}
		2. 执行Chain
		Chain(annotate("first"),annotate("second")):
			2.1 second
			next:= func(ctx context.Context, request interface{}) (interface{}, error) {
						fmt.Println("second", "pre")
						return next(ctx, request)
						fmt.Println("second", "post")
					}

			2.2 first
			next:=func(ctx context.Context, request interface{}) (interface{}, error) {
						fmt.Println("first", "pre")
						return func(ctx context.Context, request interface{}) (interface{}, error) {
								fmt.Println("second", "pre")
								return next(ctx, request)
								fmt.Println("second", "post")
							}(ctx, request)
						fmt.Println("first", "post")
					}

			2.3 返回
			next:= func(next Endpoint) Endpoint {
						func(ctx context.Context, request interface{}) (interface{}, error) {
							fmt.Println("first", "pre")
							return func(ctx context.Context, request interface{}) (interface{}, error) {
									fmt.Println("second", "pre")
									return next(ctx, request)
									fmt.Println("second", "post")
								}(ctx, request)
							fmt.Println("first", "post")
						}
			}

		3.Chain(
			annotate("first"),
			annotate("second"),
		)(myEndpoint):
				 func(next Endpoint) Endpoint {
						func(ctx context.Context, request interface{}) (interface{}, error) {
							fmt.Println("first", "pre")
							return func(ctx context.Context, request interface{}) (interface{}, error) {
									fmt.Println("second", "pre")
									return myEndpoint(ctx, request)
									fmt.Println("second", "post")
								}(ctx, request)
							fmt.Println("first", "post")
						}
				}(myEndpoint)

			4.Chain(
			annotate("first"),
			annotate("second"),
		)(myEndpoint)(context.Background(), struct{}{}):
				func(ctx context.Context, request interface{}) (interface{}, error) {
							fmt.Println("first", "pre")
							return func(ctx context.Context, request interface{}) (interface{}, error) {
									fmt.Println("second", "pre")
									return myEndpoint(ctx, request)
									fmt.Println("second", "post")
								}(ctx, request)
							fmt.Println("first", "post")
				}(context.Background(), struct{}{})
	*/
}
func TestChainMerge(t *testing.T) {
	a := Chain(
		annotate("first"),
		annotate("second"),
		annotate("third"),
	)
	b := ChainMerge(a,
		annotate("four"),
		annotate("five"),
		annotate("six"),
	)(myEndpoint)
	if _, err := b(context.Background(), struct{}{}); err != nil {
		panic(err)
	}

}

func annotate(s string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			fmt.Println(s, "pre")
			defer fmt.Println(s, "post")
			return next(ctx, request)
		}
	}
}

func myEndpoint(context.Context, interface{}) (interface{}, error) {
	fmt.Println("my endpoint!")
	return struct{}{}, nil
}
