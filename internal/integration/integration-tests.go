package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const url string = "http://localhost:8080/v1/posts"

var (
	title   string = "IntegrationTest"
	content string = "<h1>Integration Test</h1>"
	tags           = []string{"int", "test", "tag"}
)

func main() {
	p, err := CreatePost(url, title, content, tags)
	if err != nil {
		panic(err)
	}

	p2, err := ReadPost(url, p.Slug)
	if err != nil {
		panic(err)
	}

	if p.ID != p2.ID {
		err := errors.New("created ID doesnt match read ID")
		panic(err)
	}

	p3, err := UpdatePost(url, p.Slug)
	if err != nil {
		panic(err)
	}

	p4, err := DeletePost(url, p3.Slug)
	if err != nil {
		panic(err)
	}

	if p4.Response {
		fmt.Println("Integration test was successful!")
	}
}

func CreatePost(url, title, content string, tags []string) (*PostResponseBody, error) {
	post := Post{title, tags, content}
	reqBody, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(reqBody))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println("response = " + string(body))

	var respBody PostResponseBody
	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func ReadPost(url, slug string) (*PostResponseBody, error) {
	resourceUrl := url + "/" + slug
	resp, err := http.Get(resourceUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("response = " + string(body))

	var respBody PostResponseBody
	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func DeletePost(url, slug string) (*DeletePostResponse, error) {
	resourceUrl := url + "/" + slug
	cli := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, resourceUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("response = " + string(body))

	var delBody DeletePostResponse
	if err := json.Unmarshal(body, &delBody); err != nil {
		return nil, err
	}

	return &delBody, nil
}

func UpdatePost(url, slug string) (*PostResponseBody, error) {
	post := Post{title + "update", append(tags, "update"), content + "update"}
	updateReq := UpdatePostRequest{Post: post}

	reqBody, err := json.Marshal(updateReq)
	if err != nil {
		return nil, err
	}

	fmt.Println("json= " + string(reqBody))

	resourceUrl := url + "/" + slug
	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPut, resourceUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("response = " + string(body))

	var respBody PostResponseBody
	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

type DeletePostResponse struct {
	Response bool `json:"response"`
}

type UpdatePostRequest struct {
	Post Post `json:"post"`
}

type Post struct {
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
	Content string   `json:"content"`
}

type PostDeleteResponseBody struct {
	ID string `json:"id"`
}

type PostResponseBody struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Content   string   `json:"content"`
	Slug      string   `json:"slug"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
