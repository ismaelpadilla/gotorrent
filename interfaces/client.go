package interfaces

type Client interface {
	Search(query string) []Torrent
	NavigateTo(torrent Torrent)
	FetchTorrentDescription(torrent Torrent) string
	FetchTorrentFiles(torrent Torrent) []TorrentFile
}
