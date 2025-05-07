package main

const (
	MAXTIMEREQ = 20
)

func (r *Router) router() {
	r.mappingRouterAdmin()
	r.mappingRouterCustomer()
}
