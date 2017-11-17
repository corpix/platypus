package scheme

const (
	Schemes = map[string]Scheme{
		Plaintext.Scheme: Plaintext,
		Secure.Scheme:    Secure,
	}
)

func (s Schemes) Get(scheme string) (Scheme, error) {
	var (
		s  Scheme
		ok bool
	)

	s, ok = Schemes[scheme]
	if !ok {
		return s, NewErrUnsupportedScheme(scheme)
	}

	return s, nil
}
