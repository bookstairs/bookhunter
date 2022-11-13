package fetcher

type telegramService struct {
	config *Config
}

func newTelegramService(config *Config) (service, error) {
	return &telegramService{
		config: config,
	}, nil
}

func (s *telegramService) size() (int64, error) {
	panic("implement me")
}

func (s *telegramService) formats(id int64) (map[Format]string, error) {
	panic("implement me")
}

func (s *telegramService) fetch(id int64, format Format, url string) (*fetch, error) {
	panic("implement me")
}
