package main

import (
	"fmt"

	"github.com/elastifeed/es-pusher/pkg/document"
)

func main() {

	d := document.Document{
		Caption:    "TestCaption",
		Content:    "TestContent",
		Url:        "http://test/testing",
		IsFromFeed: true,
		FeedUrl:    "http://rssfeedgen.com/test.xml",
	}

	res, _ := d.Dump()

	fmt.Println(string(res))

}
