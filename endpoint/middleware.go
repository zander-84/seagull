package endpoint

// Middleware is a chainable behavior modifier for endpoints.
type Middleware func(next HandlerFunc) HandlerFunc

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(m ...Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

func ChainMerge(outer Middleware, m ...Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}

		if outer == nil {
			return next
		} else {
			return outer(next)
		}
	}
}
