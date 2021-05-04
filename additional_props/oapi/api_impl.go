package oapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// BlogAPIImpl implements OpenAPI-based endpoints.
type BlogAPIImpl struct{
	posts map[int]Post
	lock sync.Mutex
	nextID int
}

func NewBlogAPIImpl() *BlogAPIImpl {
	return &BlogAPIImpl{
		posts: make(map[int]Post),
		nextID: 1,
	}
}

func (b *BlogAPIImpl) FindPosts(w http.ResponseWriter, r *http.Request) {
	result := make([]Post, 0)
	for _, v := range b.posts {
		result = append(result, v)
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func (b *BlogAPIImpl) AddPost(w http.ResponseWriter, r *http.Request) {
	var newPost NewPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		log.Println(err)
		b.sendBlogAPIError(w, http.StatusBadRequest, "Invalid format for NewPost")
		return
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	fmt.Printf("additionalProperties: %#v\n", newPost.AdditionalProperties)

	var post Post
	post.Title = newPost.Title
	post.Content = newPost.Content
	post.Id = b.nextID
	b.nextID++

	b.posts[post.Id] = post

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

func (b *BlogAPIImpl) sendBlogAPIError(w http.ResponseWriter, code int, message string) {
	err := Error{
		Code:    code,
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(err)
}

