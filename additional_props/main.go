package main

import (
	"encoding/json"
	"fmt"
	"github.com/bellwood4486/oapi-codegen-samples/additional_props/oapi"
	"net/http"
	"sync"
)

func sendBlogAPIError(w http.ResponseWriter, code int, message string) {
	err := oapi.Error{
		Code:    code,
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(err)
}

// BlogAPIImpl implements OpenAPI-based endpoints.
type BlogAPIImpl struct{
	posts map[int]oapi.Post
	lock sync.Mutex
	nextID int
}

func NewBlogAPIImpl() *BlogAPIImpl {
	return &BlogAPIImpl{
		posts: make(map[int]oapi.Post),
	}
}

func (b *BlogAPIImpl) FindPosts(w http.ResponseWriter, r *http.Request) {
	result := make([]oapi.Post, 0)
	for _, v := range b.posts {
		result = append(result, v)
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func (b *BlogAPIImpl) AddPost(w http.ResponseWriter, r *http.Request) {
	var newPost oapi.NewPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		sendBlogAPIError(w, http.StatusBadRequest, "Invalid format for NewPost")
		return
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	var post oapi.Post
	post.Title = newPost.Title
	post.Content = newPost.Content
	post.Id = b.nextID
	b.nextID++

	b.posts[post.Id] = post

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

func main() {
	r := oapi.Handler(NewBlogAPIImpl())
	addr := ":8000"
	fmt.Printf("listen %s ...\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
