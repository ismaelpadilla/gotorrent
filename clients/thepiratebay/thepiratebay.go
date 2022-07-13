package thepiratebay

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ismaelpadilla/gotorrent/interfaces"
)

type pirateBayTorrent struct {
	Name     string
	InfoHash string `json:"info_hash"`
	Leechers string
	Seeders  string
	Size     string
	Added    string
}

type pirateBay struct{}

func New() interfaces.Client {
	return pirateBay{}
}

func (p pirateBay) HealthCheck() bool {
	return true
}

func (p pirateBay) Search(a string) []interfaces.Torrent {
	url := "https://apibay.org/q.php?q=" + a
	result, _ := http.Get(url)

	body, _ := io.ReadAll(result.Body)
	var bodyParsed []pirateBayTorrent
	json.Unmarshal(body, &bodyParsed)

	torrents := make([]interfaces.Torrent, len(bodyParsed))
	for i, pbt := range bodyParsed {
		torrents[i] = pbt.convert()
	}
	return torrents
}

//
func (p pirateBayTorrent) convert() interfaces.Torrent {
	magnetLink := "magnet:?xt=urn:btih:" + p.InfoHash
	size, _ := strconv.Atoi(p.Size)
	seeders, _ := strconv.Atoi(p.Seeders)
	leechers, _ := strconv.Atoi(p.Leechers)

	return interfaces.Torrent{
		Title:      p.Name,
		MagnetLink: magnetLink,
		Size:       size,
		Uploaded:   p.Added,
		Seeders:    seeders,
		Leechers:   leechers,
	}
}
