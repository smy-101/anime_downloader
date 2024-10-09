package main

import (
	"github.com/smy-101/anime_downloader/cmd"
	"github.com/smy-101/anime_downloader/config"
)

func main() {
	config.InitConfig()
	cmd.Execute()
}
