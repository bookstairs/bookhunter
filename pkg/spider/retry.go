package spider

func (c *Client) Retry(f func() error) (err error) {
	for i := 0; i < c.config.Retry; i++ {
		if err = f(); err == nil {
			break
		}
	}

	return
}
