package discovery

// PlainResolver returns the given hostname as it is.
// It might be useful to resolve services using the default DNS lookup mechanism.
type PlainResolver struct {
	address string
}

func (d PlainResolver) Resolve() (string, bool) {
	return d.address, true
}

func NewPlainResolver(address string) PlainResolver {
	return PlainResolver{address}
}
