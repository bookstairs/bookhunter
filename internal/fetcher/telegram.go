package fetcher

import "github.com/bookstairs/bookhunter/internal/driver"

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

func (s *telegramService) formats(id int64) (map[Format]driver.Share, error) {
	panic("implement me")
}

func (s *telegramService) fetch(id int64, format Format, share driver.Share) (*fetch, error) {
	panic("implement me")
}
