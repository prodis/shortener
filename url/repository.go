package url

type Repository interface {
	Exists(id string) bool
	FindById(id string) *Url
	FindByUrl(url string) *Url
	Save(url Url) error
}
