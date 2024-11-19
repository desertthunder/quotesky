// Request types and handlers
package api

import (
	"fmt"
	"time"
)

const CreateRecord string = "com.atproto.repo.createRecord"
const PostType string = "app.bsky.feed.post"

type Message struct {
	Content  string
	Hashtags []string
}

func (m Message) Format() string {
	return fmt.Sprintf(
		"%s\n%s\n",
		m.Content,
		time.Now().Format(time.Kitchen),
	)
}

type PostRecord struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

func (p PostRecord) CreatedAtTime() (time.Time, error) {
	return time.Parse(time.RFC3339, p.CreatedAt)
}

// Request Body for Post Creation
//
// https://docs.bsky.app/docs/advanced-guides/posts
type PostRequest struct {
	// Session did
	Repo       string     `json:"repo"`
	Collection string     `json:"collection"`
	Record     PostRecord `json:"record"`
}

func BuildPost(m Message) *PostRecord {
	return &PostRecord{
		Type:      PostType,
		Text:      m.Format(),
		CreatedAt: time.Now().Format("2006-01-02T15:04:05.000000Z"),
	}
}

func BuildPostRequest(r string, c string, p PostRecord) *PostRequest {
	return &PostRequest{r, c, p}
}
