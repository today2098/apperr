package apperr

type Body interface {
	Is(target Body) bool
	Copy() Body
}
