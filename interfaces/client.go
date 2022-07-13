package interfaces

type Client interface {
	HealthCheck() bool
	Search(query string) []Torrent
}
