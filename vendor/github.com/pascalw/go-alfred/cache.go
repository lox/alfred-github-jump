package alfred

import (
	"os"
	"encoding/json"
	"time"
)

func WriteCache(obj interface{}) {
	cacheFile, err := os.Create("cache")
	if err != nil { panic(err) }

	defer cacheFile.Close()
	json.NewEncoder(cacheFile).Encode(obj)
}

func ReadCache(into interface {}, expireAfter time.Duration) ( err error) {
	cacheFile, err := os.Open("cache")
	if err != nil { return err }

	defer cacheFile.Close()

	fileInfo, _ := cacheFile.Stat()
	if time.Since(fileInfo.ModTime()) > expireAfter {
		return nil
	}

	json.NewDecoder(cacheFile).Decode(&into)
	return nil
}
