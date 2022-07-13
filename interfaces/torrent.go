package interfaces

import "github.com/inhies/go-bytesize"

type Torrent struct {
	Title      string
	MagnetLink string
	Size       int
	Uploaded   string
	Seeders    int
	Leechers   int
}

func (t Torrent) GetPrettySize() string {
	return bytesize.New(float64(t.Size)).String()
}
