package closer

type Closer struct {
	fns []func() error
}

func (c *Closer) Close() error {
	var retErr error
	for _, fn := range c.fns {
		if err := fn(); err != nil {
			retErr = err
		}
	}
	return retErr
}

func (c *Closer) OnClose(fn func() error) {
	c.fns = append(c.fns, fn)
}

func (c *Closer) Reset() {
	if c.fns != nil {
		c.fns = c.fns[:0]
	}
}
