package post

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

const accessToken = "EAACyW7TzY6kBALaMI39X6E55xRfx6lbnIqPBpYqU28tNn5ZBjXLB8n4sxvhb8bXdwOqJMQTO4dxXzn82H7klnDDXlYAkzHjw5pPjBatSwwVUZAVVXZAf8wBRTZA4v4O9s1bZBw6DBQhXbnBgSj92ZAJIQYGdS5sg8ZD"

// Post sruct
type Post struct {
	SharesCount      uint16
	LikesCount       uint16
	CommentsCount    uint16
	ImageAttachments *[]Image
	Message          string
	Link             string
	ID               string
	Caption          string
	Descripton       string
	CoverImage       *Image
	CreatedAt        time.Time
}

// Image struct. will be used for Post as Attachments
type Image struct {
	Height, Width uint8
	URL           string
}

// Posts array
type Posts []Post

func parsePostJSON(value gjson.Result) *Post {
	post := new(Post)
	post.SharesCount = uint16(value.Get("shares.count").Uint())
	post.LikesCount = uint16(value.Get("likes.data.summary.total_count").Uint())
	post.CommentsCount = uint16(value.Get("comments.data.summary.total_count").Uint())
	post.Link = value.Get("link").String()
	post.Caption = value.Get("caption").String()
	post.Descripton = value.Get("description").String()
	post.Message = value.Get("message").String()
	post.ID = value.Get("id").String()
	post.CreatedAt = time.Unix(value.Get("created_time").Int(), 0)
	post.CoverImage = &Image{
		URL:    value.Get("attachments.data.0.media.image.src").String(),
		Width:  uint8(value.Get("attachments.data.0.media.image.width").Uint()),
		Height: uint8(value.Get("attachments.data.0.media.image.height").Uint())}
	images := make([]Image, 0, 20)
	value.Get("attachments.data.subattachments.data").ForEach(func(_, value gjson.Result) bool {
		image := &Image{
			URL:    value.Get("media.image.src").String(),
			Width:  uint8(value.Get("media.image.width").Uint()),
			Height: uint8(value.Get("media.image.height").Uint())}
		images = append(images, *image)
		return value.Exists()
	})

	post.ImageAttachments = &images
	return post

}

// GetPagePosts function
func GetPagePosts(pageID string, limit int) *Posts {

	requestURL, _ := url.Parse(fmt.Sprintf("https://graph.facebook.com/v2.11/%s/posts", pageID))
	query := requestURL.Query()
	query.Set("fields", "shares,likes.summary(true).limit(0),comments.summary(true).limit(0),message,attachments{subattachments,media},link,caption,created_time,description")
	query.Set("limit", strconv.Itoa(limit))
	query.Set("access_token", accessToken)
	query.Set("date_format", "U")
	requestURL.RawQuery = query.Encode()
	response, err := http.Get(requestURL.String())

	if err != nil {
		log.Fatal("Failed to get data", err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Failed to get response body", err)
	}

	posts := make(Posts, 0, limit)

	gjson.Get(string(body), "data").ForEach(func(_, value gjson.Result) bool {
		post := parsePostJSON(value)

		posts = append(posts, *post)

		return value.Exists()
	})

	return &posts
}
