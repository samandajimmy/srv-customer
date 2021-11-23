package nhttp

import "net/http"

type HandlerChain []http.Handler

func (hc HandlerChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range hc {
		// If request is ending, then break the chain
		rx := NewRequest(r)
		if rx.HasEnded() {
			return
		}
		h.ServeHTTP(w, r)
	}
}

type HandlerChainBuilder struct {
	chain []http.Handler
}

func NewHandlerChain(handlers ...http.Handler) HandlerChainBuilder {
	if len(handlers) == 0 {
		handlers = make([]http.Handler, 0)
	}
	return HandlerChainBuilder{chain: handlers}
}

func (b *HandlerChainBuilder) Next(h0 http.Handler, hn ...http.Handler) *HandlerChainBuilder {
	b.chain = append(b.chain, h0)
	if len(hn) > 0 {
		b.chain = append(b.chain, hn...)
	}
	return b
}

func (b *HandlerChainBuilder) Build() HandlerChain {
	return b.chain
}
