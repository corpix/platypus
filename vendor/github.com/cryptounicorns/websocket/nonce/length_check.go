package nonce

func LengthCheck(n Nonce, s int) error {
	var (
		l = base64DecodedLen(n)
	)

	if l != s {
		return NewErrNonceLengthMismatch(s, l)
	}

	return nil
}
