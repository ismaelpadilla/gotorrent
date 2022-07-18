package interfaces

import "github.com/inhies/go-bytesize"

type Torrent struct {
	Client      Client
	ID          string
	Title       string
	Description string
	Files       []TorrentFile
	MagnetLink  string
	Size        int
	Uploaded    string
	Seeders     int
	Leechers    int
}

func (t Torrent) GetPrettySize() string {
	return bytesize.New(float64(t.Size)).String()
}

func (t Torrent) FetchDescription() string {
	return t.Client.FetchTorrentDescription(t)
}

func (t Torrent) FetchFiles() []TorrentFile {
	return t.Client.FetchTorrentFiles(t)
}
