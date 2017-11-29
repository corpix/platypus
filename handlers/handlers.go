package handlers

type Handlers []Handler

func (hs Handlers) Close() error {
	var (
		err error
	)

	for _, h := range hs {
		err = h.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
