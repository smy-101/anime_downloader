package torrentmanager

type Torrent struct {
	URL        string
	Title      string
	Content    []byte
	Downloaded bool
}
