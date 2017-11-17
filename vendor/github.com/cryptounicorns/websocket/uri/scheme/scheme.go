package scheme

type Scheme struct {
	Scheme string
	Port   string
}

func (s Scheme) String() string {
	return s.Scheme
}
