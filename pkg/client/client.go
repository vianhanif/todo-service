package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vianhanif/go-pkg/httpclient"
)

// Client .
type Client interface {
	Create(ctx context.Context, payload *TodoPayload) (*Todo, error)
	List(ctx context.Context, query *Query) ([]*Todo, error)
}

// Query .
type Query struct {
	Filters []Filter
}

// Filter .
type Filter struct {
	Key   string
	Value string
}

// AddFilter .
func (c *Query) AddFilter(key string, value string) {
	c.Filters = append(c.Filters, Filter{Key: key, Value: value})
}

// TodoPayload .
type TodoPayload struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Todo .
type Todo struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Detail    string     `json:"detail"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

const (
	apiURL     = ""
	apiTimeout = 80 * time.Second
)

// Config .
type Config struct {
	APIURL            string
	HTTPClient        *http.Client
	Token             string
	MaxNetworkRetries int // retrying count defaults to 1
}

// App .
type App struct {
	Token string
	Cli   *httpclient.Client
}

// Create .
func (c *App) Create(ctx context.Context, payload *TodoPayload) (*Todo, error) {
	path := "private/v1/todos"
	jsonBody, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.Cli.APIURL, path), bytes.NewBuffer(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	if err != nil {
		return nil, err
	}

	todo := &Todo{}
	resp, err := c.Cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errDecode := c.Cli.BuildResponse(path, resp.StatusCode, jsonBody, resBody, &todo)
	if errDecode != nil {
		return nil, errDecode
	}
	return todo, nil
}

// List .
func (c *App) List(ctx context.Context, query *Query) ([]*Todo, error) {
	path := "private/v1/todos%s"

	var params string

	if query != nil {
		params += "?"
	}

	for _, filter := range query.Filters {
		params += fmt.Sprintf("%s=%s", filter.Key, filter.Value)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(path, c.Cli.APIURL, params), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	if err != nil {
		return nil, err
	}

	resp, err := c.Cli.Do(req)
	if err != nil {
		return nil, err
	}

	list := []*Todo{}

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	errDecode := c.Cli.BuildResponse(path, resp.StatusCode, nil, resBody, &list)
	if errDecode != nil {
		return nil, errDecode
	}

	return list, nil
}

// NewClient .
func NewClient(config *Config) Client {
	cli := httpclient.NewHTTPClient(httpclient.Config{
		APIURL:            config.APIURL,
		MaxNetworkRetries: config.MaxNetworkRetries,
	})
	return &App{
		Token: config.Token,
		Cli:   cli,
	}
}
