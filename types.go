package main

import (
	"errors"
	"time"
)

var (
	errNotFound = errors.New("sunset not found")
)

type sunsetResult struct {
	Sunset    string    `json:"sunset"`
	Timestamp time.Time `json:"timestamp"`
}

type sunsetFinder interface {
	Query(location string) (sunsetResult, error)
}
