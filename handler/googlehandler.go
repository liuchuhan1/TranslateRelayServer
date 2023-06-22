package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type PostJsonBody struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	Target string `json:"target"`
	Source string `json:"source"`
}

type ResultBody struct {
	Text         string `json:"text"`
	From         string `json:"from"`
	To           string `json:"to"`
	Result       string `json:"result"`
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

var (
	ServerSelect   = 1
	ServerIndexMap = map[int]string{
		1: "https://translate.google.com/translate_a/single?client=gtx&sl=",
		2: "https://translate.amz.wang/translate_a/single?client=gtx&sl=",
	}
)

func GoogleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getQuery := r.URL.Query()
		getStatus := getQuery.Get("status")
		getServerIndex, err := strconv.Atoi(getStatus)
		if err != nil {
			fmt.Print("status code is not a number")
			io.WriteString(w, "Parameter Invalid")
			return
		}
		if _, exist := ServerIndexMap[getServerIndex]; !exist {
			ServerSelect = 1
		} else {
			ServerSelect = getServerIndex
		}
		isokornot, ServerSelect := ChooseServer(ServerSelect)
		if isokornot {
			os.WriteFile("ServerSel.txt", []byte(string(ServerSelect)), fs.ModeType)
			io.WriteString(w, "InitSuccess")
			return
		} else {
			io.WriteString(w, "InitFailed")
			return
		}
	}
	if r.Method == "POST" {
		pjb := PostJsonBody{}
		rbody, _ := io.ReadAll(r.Body)
		trbody := string(rbody)
		json.Unmarshal([]byte(trbody), &pjb)
		w.Header().Set("content-type", "text/json")
		//file, err := os.Open("ServerSel.txt")
		//if err != nil {
		//	ServerSelect = 2
		//} else {
		//	var s = make([]byte, 4)
		//	file.Read(s)
		//	ServerSelect, _ = strconv.Atoi(string(s))
		//}
		fullurl := ServerIndexMap[ServerSelect] + pjb.Source + "&tl=" + pjb.Target + "&dt=t&q=" + url.QueryEscape(pjb.Text)
		resp, err := http.Get(fullurl)
		if err != nil {
			result, _ := json.Marshal(ResultBody{
				Text:         "Remote Translate Server Cant Response",
				From:         pjb.Source,
				To:           pjb.Target,
				Result:       "Error",
				ErrorCode:    "1",
				ErrorMessage: "Remote Translate Server Cant Response",
			})
			io.WriteString(w, string(result))
			return
		} else {
			if resp.StatusCode != 200 {
				result, _ := json.Marshal(ResultBody{
					Text:         "Remote Translate Server Return StatusCode:" + string(resp.StatusCode),
					From:         pjb.Source,
					To:           pjb.Target,
					Result:       "Error",
					ErrorCode:    "1",
					ErrorMessage: "Remote Translate Server Return StatusCode:" + string(resp.StatusCode),
				})
				io.WriteString(w, string(result))
				return
			}
			respbody, _ := io.ReadAll(resp.Body)
			var respInterface []interface{}
			json.Unmarshal(respbody, &respInterface)
			a, _ := respInterface[0].([]interface{})
			var d string
			i := 0
			for i < len(a) {
				b, _ := a[i].([]interface{})
				c, _ := b[0].(string)
				d += c
				i++
				fmt.Print(c)
			}

			result, _ := json.Marshal(ResultBody{
				Text:         pjb.Text,
				From:         pjb.Source,
				To:           pjb.Target,
				Result:       d,
				ErrorCode:    "0",
				ErrorMessage: "",
			})
			w.Write(result)
			return
		}
	}
}

func ChooseServer(index int) (bool, int) {
	for url, exist1 := ServerIndexMap[index]; exist1; {
		InitOK := DoGet(url + "zh-CN&tl=en&dt=t&q=%E4%BD%A0%E5%A5%BD")
		if InitOK {
			return true, index
		} else {
			index += 1
		}
	}
	return false, 0
}
