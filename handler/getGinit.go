package handler

import (
	"fmt"
	"io"
	"net/http"
)

func DoGet(url string) bool {
	{
		resp, err := http.Get(url)
		if err != nil {
			return false
		} else {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return false
			} else {
				jsonStr := string(body)
				fmt.Print(jsonStr)
				strhello := jsonStr[4:9]
				fmt.Print(strhello)
				if strhello == "Hello" {
					return true
				} else {
					return false
				}
			}
		}
	}
}
