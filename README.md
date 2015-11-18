# gotube
##Description
Gotube is a YouTube downloader using go language.
go language is a new light-weight language developed by Google, 
it provides not only many powerful libraries, but also a simple multi-threading syntax.

##Installation
- Install go from [https://golang.org/](https://golang.org/)
- Set up go environment as in [https://golang.org/doc/install](https://golang.org/doc/install)
  - (For Un*x users) ```export GOPATH=$HOME/path/to/workspace/```
  - (For Windows users) Add GOPATH=path/to/workspace/ in environment variables
- For command-line usage
  - Download the zip file or clone the project to path/to/workspace/src/
  - Type the following in command/terminal
    - ```go install gotube```
    - ```gotube -h```
- For library usage
  - You don't need to do anything

##Command-line usage
You can check the video list from a url or video id:
```
gotube -l -url https://www.youtube.com/watch?v=C0DPdy98e4c
```
Filter the video list by specifying video quality:
```
gotube -l -id C0DPdy98e4c -q medium
```
Download the video with default resolution:
```
gotube -d -url https://www.youtube.com/watch?v=C0DPdy98e4c
```
Download the video with specified resolution:
```
gotube -d -id C0DPdy98e4c -q medium -ext video/mp4
```
Specify a video repository using relative path:
```
gotube -d -id C0DPdy98e4c -q medium -ext video/mp4 -rep ./videos
```

#Library usage
```go
package main

import (
  . "github.com/KeluDiao/gotube/api"
  "log"
)

func main() {
	idList := [...]string{"shLTrG_noKo", "Ojv7tKpzkyM", "GahnMbhmt7g"}
	cl := Client{VideoRepository: "Curry_highlights"}
	for _, id := range idList {
		vl, err := cl.GetVideoListFromId(id)
		if err != nil {
			log.Fatal(err)
		}
		err = vl.Download(cl, "", "video/mp4")
		if err != nil {
			log.Fatal(err)
		}
	}
}
```

This program is still under-developing. More interesting functionalities will be added into it soon! 
