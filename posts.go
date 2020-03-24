package main

import (
	"context"
	"errors"
	"github.com/gomarkdown/markdown"
	slugs "github.com/gosimple/slug"
	"github.com/zhughes3/grpc/models"
	"strings"
)

func (s *server) ReadPost(ctx context.Context, p *models.ReadPostRequest) (*models.Post, error) {
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

	id, err := s.db.insertPost(title, slug, html)
	if err != nil {
		return nil, err
	}

	for _, t := range tags {
		_, err := s.db.insertTag(t, id)
		if err != nil {
			return nil, err
		}
	}

	return &models.Post{
		Title: title,
		Tags: tags,
		Content: html,
		Slug: slug,
	}, nil
}
