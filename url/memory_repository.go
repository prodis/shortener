package url

type memoryRepository struct {
	urls map[string]*Url
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{make(map[string]*Url)}
}

func (r *memoryRepository) Exists(id string) bool {
	_, exist := r.urls[id]
	return exist
}

func (r *memoryRepository) FindById(id string) *Url {
	return r.urls[id]
}

func (r *memoryRepository) FindByUrl(url string) *Url {
	for _, u := range r.urls {
		if u.Target == url {
			return u
		}
	}

	return nil
}

func (r *memoryRepository) Save(url Url) error {
	r.urls[url.Id] = &url
	return nil
}
