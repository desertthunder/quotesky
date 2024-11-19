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
		m.Hashtags,
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

func (p *PostRecord) SetCreatedAt() {
	t := time.Now()

	p.CreatedAt = t.Format(time.RFC3339)
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

func BuildPost(m Message, r string) *PostRecord {
	p := PostRecord{Type: PostType, Text: m.Format()}
	p.SetCreatedAt()

	return &p
}
