package main

import (
	"fmt"

	"github.com/ismaelpadilla/gotorrent/clients/thepiratebay"
)

func main() {
	fmt.Println("gotorrent")

	client := thepiratebay.New()

	torrents := client.Search("aa")
	fmt.Println(torrents[0:2])
}
