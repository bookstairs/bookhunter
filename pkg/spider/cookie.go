package spider

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sync"

	"github.com/bibliolater/bookhunter/pkg/log"
)

type jar struct {
	mu      sync.Mutex // mu locks the remaining fields.
	cookies map[string][]*http.Cookie
	subJar  http.CookieJar // The backend cookie jar.
	path    string         // The path for saving cookies.
}

func (j *jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.mu.Lock()
	defer j.mu.Unlock()

	s := u.String()
	if c, ok := j.cookies[s]; ok {
		c = append(c, cookies...)
		j.cookies[s] = c
	} else {
		j.cookies[s] = cookies
	}

	// Create or open the file.
	file, err := os.OpenFile(j.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&j.cookies); err != nil {
		panic(err)
	}

	j.subJar.SetCookies(u, cookies)
}

func (j *jar) Cookies(u *url.URL) []*http.Cookie {
	return j.subJar.Cookies(u)
}

// loadCookies would load all the cookies from file.
func (j *jar) loadCookies() error {
	if _, err := os.Stat(j.path); err == nil {
		log.Info("Found cookie file, load it.")

		// Open file.
		file, err := os.Open(j.path)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		// Loading cookies.
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&j.cookies); err != nil {
			return err
		}

		// Set cookies.
		for k, v := range j.cookies {
			if u, err := url.Parse(k); err != nil {
				return err
			} else {
				j.subJar.SetCookies(u, v)
			}
		}
	}

	return nil
}

// NewCookieJar would create a cookie jar instance. It supports auto saving policy.
func NewCookieJar(path string) (http.CookieJar, error) {
	subJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	j := &jar{
		subJar:  subJar,
		path:    path,
		cookies: make(map[string][]*http.Cookie),
	}

	if err := j.loadCookies(); err != nil {
		return nil, err
	}

	return j, nil
}
