package apperr

// Body in an interface that represents any information about response.
type Body interface {
	// Is reports whether the Body matches target (api called in apperr.Error.Is function).
	Is(target Body) bool
	// Clone returns the clone (api used in apperr.Error.Wrap funtion).
	Clone() Body
}

// H represents Body like map.
type H map[string]any

var _ (Body) = (H)(nil)

// Is reports true anytime (api called in apperr.Error.Is function).
func (h H) Is(_ Body) bool {
	return true
}

// Clone creates the clone (api used in apperr.Error.Wrap funtion).
func (h H) Clone() Body {
	clone := make(H, len(h))
	for key, elem := range h {
		clone[key] = elem
	}
	return clone
}
