package document

import (
	"fmt"
	"testing"
	"time"

	"encoding/json"
)

// Testpair containing key and corresponding value for automated testing
type testpair struct {
	key      string
	expected string
}

// Some Test values
var (
	created, _ = time.Parse(time.RFC3339, "2019-10-10T10:00:00.000Z") // Use a fixed value as time.Now() is too precise
	caption    = "Super Important Caption"
	content    = "Contant Blaaaaablabla 1234!"
	url        = "http=//test.super.important/gotestyourself.html"
	isFromFeed = true
	feedUrl    = "http://feed.wow.com/rss.xml"
)

// Generate Sample Document for testing
func genSampleDoc() Document {
	return Document{
		Created:    created,
		Caption:    caption,
		Content:    content,
		Url:        url,
		IsFromFeed: isFromFeed,
		FeedUrl:    feedUrl,
	}
}

func TestDocumentDump(t *testing.T) {
	d, _ := genSampleDoc().Dump()

	var dict map[string]interface{}

	if json.Unmarshal(d, &dict) != nil {
		t.Error("Dump generated invalid JSON")
	}

	for _, pair := range []testpair{
		{"created", created.Format(time.RFC3339)},
		{"caption", caption},
		{"content", content},
		{"url", url},
		{"feedUrl", feedUrl},
	} {
		if dict[pair.key] != pair.expected {
			t.Errorf("JSON field mismatch on key '%s' - expected: '%s', found: '%s'", pair.key, pair.expected, dict[pair.key])
		}
	}
}

func TestDocumentLoad(t *testing.T) {
	orig := genSampleDoc()
	d, _ := orig.Dump()
	loaded, _ := Load(d)

	fmt.Println(orig)
	fmt.Println(loaded)

	if orig != loaded {
		t.Error("Mismatch when loading a dumped Document")
	}
}
