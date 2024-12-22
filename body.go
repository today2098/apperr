package apperr

// Body in an interface that represents any information about response.
type Body interface {
	// Is reports whether the Body matches target (api called in apperr.Error.Is function).
	Is(target Body) bool
	// Copy returns the copy (api used in apperr.Error.Wrap funtion).
	Copy() Body
}
