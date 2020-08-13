package http

// ErrRouteConflict is returned when HTTP route conflict is detected.
type ErrRouteConflict string

func (e ErrRouteConflict) Error() string { return string(e) }
