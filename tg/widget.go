package tg

type Maker[V any] interface {
	Make(*Context) V
}


