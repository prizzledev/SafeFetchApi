package safeFetch

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls "github.com/bogdanfinn/tls-client"
)

type SafeFetch struct {
	client      tls.HttpClient
	maxAttempts int
	counter     int
}

func NewSafeFetch(max int) *SafeFetch {
	options := []tls.HttpClientOption{
		tls.WithTimeout(10),
		tls.WithClientProfile(tls.Chrome_110),
	}

	client, err := tls.NewHttpClient(nil, options...)
	if err != nil {
		panic(err)
	}

	return &SafeFetch{
		maxAttempts: max,
		client:      client,
		counter:     0,
	}
}

func (s *SafeFetch) SetMaxAttempts(maxAttempts int) {
	s.maxAttempts = maxAttempts
}

func (s *SafeFetch) generateToken(url string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://docs.google.com/gview?url=%s", url), nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"accept":             []string{"*/*"},
		"accept-language":    []string{"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"sec-ch-ua":          []string{"\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\""},
		"sec-ch-ua-mobile":   []string{"?0"},
		"sec-ch-ua-platform": []string{"\"Windows\""},
		"sec-fetch-dest":     []string{"empty"},
		"sec-fetch-mode":     []string{"cors"},
		"sec-fetch-site":     []string{"same-origin"},
		"Referer":            []string{fmt.Sprintf("https://docs.google.com/gview?url=%s", url)},
		"Referer-Policy":     []string{"strict-origin-when-cross-origin"},

		http.HeaderOrderKey: []string{
			"accept",
			"accept-language",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"Referer",
			"Referer-Policy",
		},
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Status code is not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	split1 := strings.Split(string(body), `text?id\u003d`)

	if len(split1) < 2 {
		return "", fmt.Errorf("Could not find token step 1")
	}

	split2 := strings.Split(split1[1], `\\u0026authuser`)

	if len(split2) < 1 {
		return "", fmt.Errorf("Could not find token step 2")
	}

	split3 := strings.Split(split2[0], `"`)

	if len(split3) < 1 {
		return "", fmt.Errorf("Could not find token step 3")
	}

	id := split3[0]

	return id, nil
}

func (s *SafeFetch) render(token string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://docs.google.com/viewerng/text?id=%s&page=0", token), nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"accept":             []string{"*/*"},
		"accept-language":    []string{"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"sec-ch-ua":          []string{"\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\""},
		"sec-ch-ua-mobile":   []string{"?0"},
		"sec-ch-ua-platform": []string{"\"Windows\""},
		"sec-fetch-dest":     []string{"empty"},
		"sec-fetch-mode":     []string{"cors"},
		"sec-fetch-site":     []string{"same-origin"},

		http.HeaderOrderKey: []string{
			"accept",
			"accept-language",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
		},
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Status code is not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content := strings.Split(string(body), "\n")

	if len(content) < 2 {
		return "", fmt.Errorf("Could not find content")
	}

	return content[1], nil
}
func (s *SafeFetch) Get(url string) (string, error) {
	fmt.Println(fmt.Sprintf("[%s] Generating token for %s", time.Now().Format("2006-01-02 15:04:05"), url))

	var token string

	for i := 0; i < s.maxAttempts; i++ {
		var err error

		token, err = s.generateToken(url)
		if err != nil {
			if err.Error() == "Status code is not 200" {
				fmt.Println(fmt.Sprintf("[%s] Failed to generate token, retrying...", time.Now().Format("2006-01-02 15:04:05")))
				continue
			}

			return "", err
		}

		if token != "" {
			break
		}
	}

	if token == "" {
		return "", fmt.Errorf(fmt.Sprintf("[%s] Maximum token generation attempts reached", time.Now().Format("2006-01-02 15:04:05")))
	}

	render, err := s.render(token)
	if err != nil {
		return "", err
	}

	s.counter++
	fmt.Println(fmt.Sprintf("[%s][%d] Successfully rendered document", time.Now().Format("2006-01-02 15:04:05"), s.counter))

	return render, nil
}
