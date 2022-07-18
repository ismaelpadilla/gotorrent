package interfaces

import "github.com/inhies/go-bytesize"

type TorrentFile struct {
	Name string
	Size int
}

func (tf TorrentFile) GetPrettySize() string {
	return bytesize.New(float64(tf.Size)).String()
}
