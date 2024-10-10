package torrentmanager

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type testEnv struct {
	db *sql.DB
}

var env *testEnv

func TestMain(m *testing.M) {
	var err error
	env = &testEnv{}
	env.db, err = initDatabase()

	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		os.Exit(1)
	}
	cleanDatabase()
	code := m.Run()
	env.db.Close()
	os.Exit(code)
}

func TestInsertTorrent(t *testing.T) {
	runTest(t, testInsertTorrent)
}

func TestGetTorrentByURL(t *testing.T) {
	runTest(t, testGetTorrentByURL)
}

func TestUpdateTorrentDownloadStatus(t *testing.T) {
	runTest(t, testUpdateTorrentDownloadStatus)
}

func TestListAllTorrents(t *testing.T) {
	runTest(t, testListAllTorrents)
}

func testInsertTorrent(t *testing.T, db *sql.DB) {
	err := insertTorrent(db, "http://example.com/torrent1", "Torrent 1", []byte("content"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func testGetTorrentByURL(t *testing.T, db *sql.DB) {
	// 首先插入一个 torrent
	err := insertTorrent(db, "http://example.com/torrent2", "Torrent 2", []byte("content"))
	if err != nil {
		t.Fatalf("Failed to insert torrent: %v", err)
	}

	torrent, err := getTorrentByURL(db, "http://example.com/torrent2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if torrent.URL != "http://example.com/torrent2" {
		t.Errorf("Expected URL to be 'http://example.com/torrent2', got %v", torrent.URL)
	}
}

func testUpdateTorrentDownloadStatus(t *testing.T, db *sql.DB) {
	err := insertTorrent(db, "http://example.com/torrent3", "Torrent 3", []byte("content"))
	if err != nil {
		t.Fatalf("Failed to insert torrent: %v", err)
	}

	err = updateTorrentDownloadStatus(db, "http://example.com/torrent3", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	torrent, err := getTorrentByURL(db, "http://example.com/torrent3")
	if err != nil {
		t.Errorf("Failed to get torrent: %v", err)
	}
	if !torrent.Downloaded {
		t.Errorf("Expected downloaded to be true, got false")
	}
}

func testListAllTorrents(t *testing.T, db *sql.DB) {
	cleanDatabase()
	err := insertTorrent(db, "http://example.com/torrent4", "Torrent 4", []byte("content"))
	if err != nil {
		t.Fatalf("Failed to insert torrent: %v", err)
	}

	err = insertTorrent(db, "http://example.com/torrent5", "Torrent 5", []byte("content"))
	if err != nil {
		t.Fatalf("Failed to insert torrent: %v", err)
	}

	torrents, err := listAllTorrents(db)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(torrents) != 2 {
		t.Errorf("Expected 2 torrents, got %v", len(torrents))
	}
}

// 辅助函数，用于运行每个测试 将数据库连接传递给测试函数
func runTest(t *testing.T, testFunc func(*testing.T, *sql.DB)) {
	testFunc(t, env.db)
}

func cleanDatabase() {
	_, err := env.db.Exec("DELETE FROM torrents")
	if err != nil {
		log.Fatalf("Failed to clean database: %v", err)
	}
}
