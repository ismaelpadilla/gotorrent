package interfaces

type Client interface {
	Search(query string) []Torrent
	FetchTorrentDescription(torrent Torrent) string
}
