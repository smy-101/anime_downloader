package torrent

import (
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
)

type Client struct {
	client    *torrent.Client
	outputDir string
}

func NewClient(outputDir string) (*Client, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = outputDir
	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{client: client, outputDir: outputDir}, nil
}

func (c *Client) DownloadFromFile(torrentPath string) error {
	t, err := c.client.AddTorrentFromFile(torrentPath)
	if err != nil {
		return err
	}
	return c.downloadTorrent(t)
}

func (c *Client) DownloadFromMagnet(magnetLink string) error {
	t, err := c.client.AddMagnet(magnetLink)
	if err != nil {
		return err
	}
	return c.downloadTorrent(t)
}

func (c *Client) downloadTorrent(t *torrent.Torrent) error {
	<-t.GotInfo()
	fmt.Printf("Torrent information retrieved: %s\n", t.Name())
	fmt.Println("Starting download...")

	t.DownloadAll()

	progressTicker := time.NewTicker(time.Second)
	defer progressTicker.Stop()

	var lastCompleted int64
	lastTime := time.Now()

	for {
		select {
		case <-progressTicker.C:
			complete := t.BytesCompleted()
			totalSize := t.Length()
			progress := float64(complete) / float64(totalSize) * 100

			// Calculate download speed
			currentTime := time.Now()
			duration := currentTime.Sub(lastTime).Seconds()
			speed := float64(complete-lastCompleted) / duration

			fmt.Printf("\rProgress: %.2f%% | Download Speed: %s/s",
				progress,
				humanize.Bytes(uint64(speed)))

			if complete == totalSize {
				fmt.Println("\nDownload completed!")
				return nil
			}

			lastCompleted = complete
			lastTime = currentTime
		}
	}
}

func (c *Client) Close() {
	c.client.Close()
}
