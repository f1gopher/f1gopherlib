// F1GopherLib - Copyright (C) 2022 f1gopher
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connection

import (
	"bufio"
	"github.com/f1gopher/f1gopherlib/f1log"
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
		cache = filepath.Join(cache, strings.Replace(url, "https://livetiming.formula1.com/static/", "", 1))
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
