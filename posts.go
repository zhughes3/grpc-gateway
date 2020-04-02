package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	slugs "github.com/gosimple/slug"
	"github.com/zhughes3/grpc/models"
)

var errEmptySlug error = errors.New("slug empty")

func (s *server) ReadPost(ctx context.Context, p *models.PostFilter) (*models.Post, error) {
	var err error
	po := &models.Post{}

	if len(p.GetSlug()) > 0 {
		po, err = s.db.readPost(p.GetSlug())
		if err != nil {
			return nil, err
		}

		return po, err
	}
	err = errors.New("slug empty")
	return nil, err
}

func (s *server) CreatePost(ctx context.Context, p *models.Post) (*models.Post, error) {
	title := p.GetTitle()
	slug := slugs.Make(title)
	tags := p.GetTags()
	content := p.GetContent()
	html := strings.ReplaceAll(string(markdown.ToHTML([]byte(content), nil, nil)), "\n", "")

	post, err := s.db.insertPost(title, slug, html)
	if err != nil {
		return nil, err
	}

	for _, t := range tags {
		_, err := s.db.insertTag(t, post.GetId())
		if err != nil {
			return nil, err
		}
	}

	return &models.Post{
		Id:        post.GetId(),
		Title:     title,
		Tags:      tags,
		Content:   html,
		Slug:      slug,
		CreatedAt: post.GetCreatedAt(),
		UpdatedAt: post.GetUpdatedAt(),
	}, nil
}

func (s *server) DeletePost(ctx context.Context, p *models.PostFilter) (*models.PostError, error) {
	var err error

	if len(p.GetSlug()) > 0 {
		_, err := s.db.deletePost(p.GetSlug())
		if err != nil {
			return nil, err
		}

		return &models.PostError{
			Response: true,
		}, err
	}
	err = errors.New("slug empty")
	return nil, err
}

func (s *server) UpdatePost(ctx context.Context, req *models.UpdatePostRequest) (*models.Post, error) {
	//var err error
	if len(req.GetSlug()) == 0 {
		fmt.Println("Slug is empty on UpdatePost request")
		return nil, errEmptySlug
	}

	updatedPost := req.GetPost()
	newSlug := slugs.Make(updatedPost.GetTitle())
	html := strings.ReplaceAll(string(markdown.ToHTML([]byte(updatedPost.GetContent()), nil, nil)), "\n", "")
	p, err := s.db.updatePost(req.GetSlug(), updatedPost.GetTitle(), newSlug, html)
	if err != nil {
		return nil, err
	}

	tags, err := s.db.updateTags(p.GetId(), req.GetPost().GetTags())
	if err != nil {
		return nil, err
	}

	p.Tags = tags

	return p, nil
}
