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

type pirateBay struct{}

func New() interfaces.Client {
	return pirateBay{}
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
		torrents[i].Client = p
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
		ID:         p.ID,
		Title:      p.Name,
		MagnetLink: magnetLink,
		Size:       size,
		Uploaded:   p.Added,
		Seeders:    seeders,
		Leechers:   leechers,
	}
}

func (p pirateBay) FetchTorrentDescription(torrent interfaces.Torrent) string {
	url := "https://apibay.org/t.php?id=" + torrent.ID
	result, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Panic(err)
	}
	var bodyParsed pirateBayTorrentDetails
	err = json.Unmarshal(body, &bodyParsed)
	if err != nil {
		log.Panic("cant unmarshall", err)
	}

	return bodyParsed.Descr
}

func (p pirateBay) FetchTorrentFiles(torrent interfaces.Torrent) []interfaces.TorrentFile {
	url := "https://apibay.org/f.php?id=" + torrent.ID
	result, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Panic(err)
	}
	var bodyParsed []pirateBayTorrentFile
	err = json.Unmarshal(body, &bodyParsed)
	if err != nil {
		log.Panic("cant unmarshall", err)
	}

	torrentFiles := make([]interfaces.TorrentFile, len(bodyParsed))
	for i, pbtf := range bodyParsed {
		torrentFiles[i] = interfaces.TorrentFile{
			Name: pbtf.Name[0],
			Size: pbtf.Size[0],
		}
	}
	return torrentFiles
}
