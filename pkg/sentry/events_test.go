package sentry_test

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/grafana/sentry-datasource/pkg/sentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseEventsQuery(t *testing.T, query string) url.Values {
	t.Helper()
	_, rawQuery, found := strings.Cut(query, "?")
	require.True(t, found, "query string %q has no query parameters", query)
	params, err := url.ParseQuery(rawQuery)
	require.NoError(t, err)
	return params
}

func TestGetEventsInput_ToQuery(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

	// Regression test for queries built outside the query editor (Scenes apps, dashboards saved
	// before 2.0, API callers): upstream 2.0.0 (#400) moved the default field list into the UI,
	// so a query with no fields sent no field= params and Sentry answered
	// "400 Bad Request No columns selected".
	t.Run("falls back to the default fields when no fields are provided", func(t *testing.T) {
		for name, fields := range map[string][]string{
			"nil fields":   nil,
			"empty fields": {},
		} {
			t.Run(name, func(t *testing.T) {
				input := sentry.GetEventsInput{
					OrganizationSlug: testOrgSlug,
					Fields:           fields,
					From:             from,
					To:               to,
					Limit:            10,
				}

				query := input.ToQuery()

				assert.True(t, strings.HasPrefix(query, "/api/0/organizations/test-org/events/?"), query)
				params := parseEventsQuery(t, query)
				assert.Equal(t, sentry.DefaultEventFields, params["field"])
			})
		}
	})

	t.Run("sends explicit fields as-is without appending defaults", func(t *testing.T) {
		input := sentry.GetEventsInput{
			OrganizationSlug: testOrgSlug,
			Fields:           []string{"title", "timestamp"},
			From:             from,
			To:               to,
			Limit:            10,
		}

		params := parseEventsQuery(t, input.ToQuery())

		assert.Equal(t, []string{"title", "timestamp"}, params["field"])
	})

	t.Run("default fields match the pre-2.0 required fields", func(t *testing.T) {
		assert.Equal(t, []string{
			"id", "title", "project", "project.id", "release",
			"count()", "epm()", "last_seen()", "level", "event.type", "platform",
		}, sentry.DefaultEventFields)
	})
}
