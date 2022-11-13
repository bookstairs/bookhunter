package fetcher

type talebookService struct {
}

func newTalebookService() (service, error) {
	return &talebookService{}, nil
}

func (t *talebookService) size() (int64, error) {
	panic("implement me")
}

func (t *talebookService) formats() ([]Format, error) {
	panic("implement me")
}

func (t *talebookService) fetch(id int64, format Format) (*fetch, error) {
	f := createFetch(nil)
	return f, nil
}
