package url

import (
	"math/rand"
	"net/url"
	"time"
)

const (
	idLength       = 5
	allowedSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-"
)

type Url struct {
	Id        string
	CreatedAt time.Time
	Target    string
}

var repo Repository

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ConfigureRepository(r Repository) {
	repo = r
}

func FindOrCreate(target string) (*Url, bool, error) {
	if url := repo.FindByUrl(target); url != nil {
		return url, false, nil
	}

	if _, err := url.ParseRequestURI(target); err != nil {
		return nil, false, err
	}

	url := Url{
		Id:        generateId(),
		CreatedAt: time.Now(),
		Target:    target,
	}
	repo.Save(url)

	return &url, true, nil
}

func Find(id string) *Url {
	return repo.FindById(id)
}

func generateId() string {
	newId := func() string {
		id := make([]byte, idLength, idLength)
		for i := range id {
			id[i] = allowedSymbols[rand.Intn(len(allowedSymbols))]
		}
		return string(id)
	}

	for {
		if id := newId(); !repo.Exists(id) {
			return id
		}
	}
}
