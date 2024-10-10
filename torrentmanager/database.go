package torrentmanager

import (
	"database/sql"
	"fmt"
)

func initDatabase() (*sql.DB, error) {
	// 打开数据库连接（如果文件不存在，SQLite 会自动创建）
	db, err := sql.Open("sqlite3", "./torrent.db")
	if err != nil {
		return nil, err
	}

	// 验证连接
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 创建 torrents 表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS torrents (
			url TEXT PRIMARY KEY,
			title TEXT,
			content BLOB,
			downloaded BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create torrents table: %v", err)
	}

	return db, nil
}

func insertTorrent(db *sql.DB, url, title string, content []byte) error {
	_, err := db.Exec("INSERT OR REPLACE INTO torrents (url, title, content) VALUES (?, ?, ?)", url, title, content)
	return err
}

func getTorrentByURL(db *sql.DB, url string) (*Torrent, error) {
	var t Torrent
	err := db.QueryRow("SELECT url, title, content, downloaded FROM torrents WHERE url = ?", url).Scan(&t.URL, &t.Title, &t.Content, &t.Downloaded)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func updateTorrentDownloadStatus(db *sql.DB, url string, downloaded bool) error {
	_, err := db.Exec("UPDATE torrents SET downloaded = ? WHERE url = ?", downloaded, url)
	return err
}

func listAllTorrents(db *sql.DB) ([]*Torrent, error) {
	rows, err := db.Query("SELECT url, title, content, downloaded FROM torrents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var torrents []*Torrent
	for rows.Next() {
		var t Torrent
		if err := rows.Scan(&t.URL, &t.Title, &t.Content, &t.Downloaded); err != nil {
			return nil, err
		}
		torrents = append(torrents, &t)
	}

	return torrents, nil
}
