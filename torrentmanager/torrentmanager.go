package torrentmanager

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
)

type TorrentManager struct {
	db *sql.DB
}

type TorrentService interface {
	AddTorrent(url, title string) error
	GetTorrent(url string) (*Torrent, error)
	UpdateDownloadStatus(url string, downloaded bool) error
	ListTorrents() ([]*Torrent, error)
	Close() error
}

func NewTorrentManager() (*TorrentManager, error) {
	db, err := initDatabase()
	if err != nil {
		return nil, err
	}

	return &TorrentManager{db: db}, nil
}

func (tm *TorrentManager) Close() error {
	return tm.db.Close()
}

func (tm *TorrentManager) AddTorrent(url, title string) error {
	// 下载 torrent 文件
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download torrent from %s: %w", url, err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read torrent content from %s: %w", url, err)
	}

	// 插入数据库
	if err := insertTorrent(tm.db, url, title, content); err != nil {
		return fmt.Errorf("failed to insert torrent into database: %w", err)
	}
	return nil
}

func (tm *TorrentManager) GetTorrent(url string) (*Torrent, error) {
	torrent, err := getTorrentByURL(tm.db, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get torrent by URL %s: %w", url, err)
	}
	return torrent, nil
}

func (tm *TorrentManager) UpdateDownloadStatus(url string, downloaded bool) error {
	if err := updateTorrentDownloadStatus(tm.db, url, downloaded); err != nil {
		return fmt.Errorf("failed to update download status for URL %s: %w", url, err)
	}
	return nil
}

func (tm *TorrentManager) ListTorrents() ([]*Torrent, error) {
	torrents, err := listAllTorrents(tm.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list torrents: %w", err)
	}
	return torrents, nil
}
