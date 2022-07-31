package thepiratebay

type pirateBay struct{}

type pirateBayTorrent struct {
	ID       string
	Name     string
	Descr    string
	InfoHash string `json:"info_hash"`
	Leechers string
	Seeders  string
	Size     string
	Added    string
}

type pirateBayTorrentDetails struct {
	Descr string
}

type pirateBayTorrentFile struct {
	Name []string
	Size []int
}
