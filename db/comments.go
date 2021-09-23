package db

import (
	// database driver
	_ "github.com/go-sql-driver/mysql"
)

// Comment contains comment data.
type Comment struct {
	Author  string
	Content string
}
