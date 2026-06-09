package url

type Repository interface {
	GetUrlByCode(short_code string) (string, error)
	AddUrl(url, key string) error
}
