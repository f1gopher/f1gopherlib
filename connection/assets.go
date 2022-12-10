package connection

import (
	"bufio"
	"f1gopherlib/f1gopherlib/f1log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AssetStore interface {
	TeamRadio(file string) ([]byte, error)
}

type assets struct {
	log   *f1log.F1GopherLibLog
	url   string
	cache string
}

func CreateAssetStore(url string, cache string, log *f1log.F1GopherLibLog) AssetStore {

	if len(cache) > 0 {
		cache = filepath.Join("./cache", strings.Replace(url, "https://livetiming.formula1.com/static/", "", 1))
	}

	return &assets{
		log:   log,
		url:   url,
		cache: cache,
	}
}

func (a *assets) TeamRadio(file string) ([]byte, error) {
	url := a.url + file

	if len(a.cache) > 0 {
		dataPath := strings.Replace(url, "https://livetiming.formula1.com/static/", "", 1)

		// If file matching url doesn't exist then retrieve
		cachedFile := filepath.Join(a.cache, dataPath)
		cachedFile, _ = filepath.Abs(cachedFile)
		f, err := os.Open(cachedFile)

		if os.IsNotExist(err) {
			f.Close()

			var resp *http.Response
			resp, err = http.Get(url)
			if err != nil {
				a.log.Errorf("Fetching team radio for '%s': %v", url, err)
				return nil, err
			}
			defer resp.Body.Close()

			scanner := bufio.NewScanner(resp.Body)

			err = os.MkdirAll(filepath.Dir(cachedFile), 0755)

			// Write body to file - using url as name
			var newFile *os.File
			newFile, err = os.Create(cachedFile)
			defer newFile.Close()
			for scanner.Scan() {
				_, err = newFile.Write(scanner.Bytes())

				// need newline for scanner to split
				newFile.WriteString("\n")
			}
			f, err = os.Open(cachedFile)
		}

		return io.ReadAll(bufio.NewReader(f))
	}

	var resp *http.Response
	resp, err := http.Get(url)
	if err != nil {
		a.log.Errorf("Fetching team radio for '%s': %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(bufio.NewReader(resp.Body))
}
