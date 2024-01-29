package authorsearch

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// loadCache loads the contents of the cache file. If it doesn't
// exist, updateCache function is called.
func (website Resource) loadCache() ([]byte, error) {
	_, err := os.Stat(website.getCacheFileName())
	if errors.Is(err, os.ErrNotExist) {
		err = website.updateCache()
	}
	if err != nil {
		return []byte{}, err
	}

	file, err := os.ReadFile(website.getCacheFileName())
	if err != nil {
		return []byte{}, err
	}
	return file, nil
}

// updateCache carries out an http get request and saves the response body
// into a file
func (website Resource) updateCache() error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(website.BaseURL + website.QueryURL)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(website.getCacheFileName(), body, 0644)
	return err
}

func (website Resource) getCacheFileName() string {
	return "cache/" + website.Name + "." + website.DataFormat
}
