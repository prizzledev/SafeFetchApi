package main

import (
	"SafeFetchApi/src/safeFetch"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	safeFetch := safeFetch.NewSafeFetch(50)

	http.HandleFunc("/api/v2/safeFetch", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			resultResponse := resultStr{
				Success: false,
				Message: err.Error(),
				Content: "",
			}

			stringRes, _ := json.Marshal(resultResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(string(stringRes)))
			return
		}

		var safeFetchRequest SafeFetchReq
		err = json.Unmarshal([]byte(body), &safeFetchRequest)
		if err != nil {
			resultResponse := resultStr{
				Success: false,
				Message: err.Error(),
				Content: "",
			}

			stringRes, _ := json.Marshal(resultResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(string(stringRes)))
			return
		}

		url := safeFetchRequest.Url

		res, err := safeFetch.Get(url)
		if err != nil {
			resultResponse := resultStr{
				Success: false,
				Message: err.Error(),
				Content: "",
			}

			stringRes, _ := json.Marshal(resultResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(string(stringRes)))
			return
		}

		var safeFetchResponse SafeFetchResp
		err = json.Unmarshal([]byte(res), &safeFetchResponse)
		if err != nil {
			resultResponse := resultStr{
				Success: false,
				Message: err.Error(),
				Content: "",
			}

			stringRes, _ := json.Marshal(resultResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(string(stringRes)))
			return
		}

		if safeFetchResponse.Mimetype == "application/json" {
			var result map[string]interface{}
			err = json.Unmarshal([]byte(safeFetchResponse.Data), &result)
			if err != nil {
				resultResponse := resultStr{
					Success: true,
					Message: "could not parse json, but send as string",
					Content: safeFetchResponse.Data,
				}

				stringRes, _ := json.Marshal(resultResponse)

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(string(stringRes)))
				return
			}

			resultResponse := resultJson{
				Success: true,
				Message: "Successfully safefetched json",
				Content: result,
			}

			stringRes, _ := json.Marshal(resultResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(string(stringRes)))
			return
		} else if safeFetchResponse.Mimetype == "text/html" {
			resultResponse := resultHtml{
				Success: true,
				Message: "",
				Content: safeFetchResponse.Data,
			}

			stringRes, _ := json.Marshal(resultResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(string(stringRes)))
			return
		} else {
			resultResponse := resultStr{
				Success: false,
				Message: "Unknown mimetype: " + safeFetchResponse.Mimetype,
				Content: "",
			}

			stringRes, _ := json.Marshal(resultResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(string(stringRes)))
			return
		}
	})

	fmt.Println(fmt.Sprintf("[%s] Listening on localhost:4501", time.Now().Format("2006-01-02 15:04:05")))
	http.ListenAndServe(":4501", nil)
}

type SafeFetchResp struct {
	Mimetype string `json:"mimetype"`
	Data     string `json:"data"`
}

type resultStr struct {
	Success bool
	Message string
	Content string
}

type resultJson struct {
	Success bool
	Message string
	Content map[string]interface{}
}

type resultHtml struct {
	Success bool
	Message string
	Content string
}

type SafeFetchReq struct {
	Url string `json:"url"`
}
