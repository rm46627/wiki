package db

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

/*
pageBucket: key - value
			ID123  |
*/

func main() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Errorf("%v", err)
	}
	defer db.Close()

}
