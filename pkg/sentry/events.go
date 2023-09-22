package sentry

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var reqFields = [...]string{
	"id",
	"title",
	"project",
	"project.id",
	"release",
	"count()",
	"epm()",
	"last_seen()",
	"failure_rate()",
	"level",
	"event.type",
	"platform",
}

type SentryEvents struct {
	Data []SentryEvent          `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

type SentryEvent struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Project         string    `json:"project"`
	ProjectId       int64     `json:"project.id"`
	Release         string    `json:"release"`
	Count           int64     `json:"count()"`
	EventsPerMinute float64   `json:"epm()"`
	LastSeen        time.Time `json:"last_seen()"`
	FailureRate     float64   `json:"failure_rate()"`
	Level           string    `json:"level"`
	EventType       string    `json:"event.type"`
	Platform        string    `json:"platform"`
}

type GetEventsInput struct {
	OrganizationSlug string
	ProjectIds       []string
	Environments     []string
	Query            string
	From             time.Time
	To               time.Time
	Sort             string
	Limit            int64
}

func (gii *GetEventsInput) ToQuery() string {
	urlPath := fmt.Sprintf("/api/0/organizations/%s/events/?", gii.OrganizationSlug)
	if gii.Limit < 1 || gii.Limit > 100 {
		gii.Limit = 100
	}
	params := url.Values{}
	params.Set("query", gii.Query)
	params.Set("start", gii.From.Format("2006-01-02T15:04:05"))
	params.Set("end", gii.To.Format("2006-01-02T15:04:05"))
	if gii.Sort != "" {
		params.Set("sort", gii.Sort)
	}
	params.Set("per_page", strconv.FormatInt(gii.Limit, 10))
	for _, field := range reqFields {
		params.Add("field", field)
	}
	for _, projectId := range gii.ProjectIds {
		params.Add("project", projectId)
	}
	for _, environment := range gii.Environments {
		params.Add("environment", environment)
	}
	return urlPath + params.Encode()
}

func (sc *SentryClient) GetEvents(gii GetEventsInput) ([]SentryEvent, string, error) {
	var out SentryEvents
	executedQueryString := gii.ToQuery()
	err := sc.Fetch(executedQueryString, &out)
	return out.Data, sc.BaseURL + executedQueryString, err
}
