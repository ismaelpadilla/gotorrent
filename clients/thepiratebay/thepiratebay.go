package thepiratebay

import (
	"encoding/json"
	"io"
	"log"
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
	result, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Panic(err)
	}
	var bodyParsed []pirateBayTorrent
	err = json.Unmarshal(body, &bodyParsed)
	if err != nil {
		log.Panic(err)
	}

	torrents := make([]interfaces.Torrent, len(bodyParsed))
	for i, pbt := range bodyParsed {
		torrents[i] = pbt.convert()
	}
	return torrents
}

//
func (p pirateBayTorrent) convert() interfaces.Torrent {
	magnetLink := "magnet:?xt=urn:btih:" + p.InfoHash
	size, err := strconv.Atoi(p.Size)
	if err != nil {
		log.Panic(err)
	}
	seeders, err := strconv.Atoi(p.Seeders)
	if err != nil {
		log.Panic(err)
	}
	leechers, err := strconv.Atoi(p.Leechers)
	if err != nil {
		log.Panic(err)
	}

	return interfaces.Torrent{
		Title:      p.Name,
		MagnetLink: magnetLink,
		Size:       size,
		Uploaded:   p.Added,
		Seeders:    seeders,
		Leechers:   leechers,
	}
}
