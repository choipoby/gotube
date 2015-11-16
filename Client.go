package main

import (
	"bytes"
	. "strings"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
YouTube client
You need to log in to view age-restricted videos
*/
type Client struct {
	userName string
	passWord string
}

/*
Request http code from url
*/
func (*Client) RequestUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
Get json data
*/
func (*Client) GetJson(httpData []byte) (map[string]interface{}, error) {
	//Find begining of json data
	var jsonBeg = "ytplayer.config = {"
	beg := bytes.Index(httpData, []byte(jsonBeg))
	if beg == -1 { //pattern not found
		return nil, PatternNotFoundError{pattern: jsonBeg}
	}
	beg += len(jsonBeg) //len(jsonBeg) returns the number of bytes in jsonBeg

	//Find offset of json data
	var unmatchedBrackets = 1
	var offset = 0
	for unmatchedBrackets > 0 {
		nextRight := bytes.Index(httpData[beg+offset:], []byte("}"))
		if nextRight == -1 {
			return nil, UnmatchedBracketsError{}
		}
		unmatchedBrackets -= 1
		unmatchedBrackets += bytes.Count(httpData[beg+offset:beg+offset+nextRight], []byte("{"))
		offset += nextRight + 1
	}

	//Load json data
	var f interface{}
	err := json.Unmarshal(httpData[beg-1:beg+offset], &f)
	if err != nil {
		return nil, err
	}
	return f.(map[string]interface{}), nil
}

/*
Get video from json data
*/
func (*Client) GetVideoList(jsonData map[string]interface{}) (videoList VideoList, err error) {
	args := jsonData["args"].(map[string]interface{})
	videoList.title = args["title"].(string)
	encodedStreamMap := args["url_encoded_fmt_stream_map"].(string)
	videoListStr := Split(encodedStreamMap, ",")
	for _, videoStr := range videoListStr {
		videoStr, err = url.QueryUnescape(videoStr)
		if err != nil {
			return
		}
		videoParams := Split(videoStr, "&")
		var video Video
		for _, param := range videoParams {
			switch {
			case HasPrefix(param, "quality"):
				video.quality = param[8:]
			case HasPrefix(param, "type"):
				video.videoType = Split(param, ";")[0][5:]
			case HasPrefix(param, "url"):
				video.url = param[4:]
			}
		}
		var missingFields []string
		if video.quality == "" {
			missingFields = append(missingFields, "quality")
		}
		if video.videoType == "" {
			missingFields = append(missingFields, "video type")
		}
		if video.url == "" {
			missingFields = append(missingFields, "url")
		}
		if missingFields != nil {
			err = MissingFieldsError{fields:missingFields}
			return
		}
		videoList.Append(video)
	}
	return
}
