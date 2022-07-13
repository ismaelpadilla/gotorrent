package interfaces

type Torrent struct {
	Title       string
	MagnetLink  string
	Size        int
	Uploaded    string
	Seeders     int
	Leechers    int
}
