package serveseq_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func RunRequests(h http.Handler, urls ...string) {
	svr := httptest.NewServer(h)
	defer svr.Close()

	for _, url := range urls {
		res, err := http.Get(svr.URL + "/" + url)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s: %s\n", url, body)
	}
}
