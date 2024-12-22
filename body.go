package apperr

// Body in an interface that represents any information about response.
type Body interface {
	// Is reports whether the Body matches target (api called in apperr.Error.Is function).
	Is(target Body) bool
	// Copy returns the copy (api used in apperr.Error.Wrap funtion).
	Copy() Body
}

// H represents Body like map.
type H map[string]any

var _ (Body) = (H)(nil)

// Is reports true anytime (api called in apperr.Error.Is function).
func (h H) Is(_ Body) bool {
	return true
}

// Copy creates the copy (api used in apperr.Error.Wrap funtion).
func (h H) Copy() Body {
	clone := make(H, len(h))
	for key, elem := range h {
		clone[key] = elem
	}
	return clone
}
