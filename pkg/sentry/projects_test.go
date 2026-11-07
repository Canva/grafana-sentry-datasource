package sentry_test

import (
	"net/http"
	"testing"

	"github.com/grafana/sentry-datasource/pkg/sentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSentryClient_GetProjects(t *testing.T) {
	t.Run("caches paginated projects under projects-pagination and hits cache on subsequent calls", func(t *testing.T) {
		pageOne := []sentry.SentryProject{
			{ID: "1", Slug: "project-1", Name: "Project One"},
		}
		pageTwo := []sentry.SentryProject{
			{ID: "2", Slug: "project-2", Name: "Project Two"},
		}

		linkHeader := `</api/0/organizations/test-org/projects/?cursor=next_cursor>; rel="next"; results="true"`
		doer := &sequencedDoer{
			responses: []*http.Response{
				createMockResponse(http.StatusOK, pageOne, map[string]string{"Link": linkHeader}),
				createMockResponse(http.StatusOK, pageTwo, nil),
			},
		}

		client, err := sentry.NewSentryClient(testBaseURL, testOrgSlug, testAuthToken, doer)
		require.NoError(t, err)

		// First call: fetches both pages
		projects, err := client.GetProjects(testOrgSlug, true, false)
		require.NoError(t, err)
		require.Len(t, projects, 2)
		assert.Equal(t, "project-1", projects[0].Slug)
		assert.Equal(t, "project-2", projects[1].Slug)
		assert.Len(t, doer.requests, 2)

		// Second call with withPagination=true: must hit cache and NOT invoke HTTP doer
		cachedProjects, err := client.GetProjects(testOrgSlug, true, false)
		require.NoError(t, err)
		require.Len(t, cachedProjects, 2)
		assert.Equal(t, "project-1", cachedProjects[0].Slug)
		assert.Equal(t, "project-2", cachedProjects[1].Slug)
		// Request count must remain 2 (cache hit)
		assert.Len(t, doer.requests, 2)
	})

	t.Run("unpaginated and paginated project caches remain isolated", func(t *testing.T) {
		unpaginatedProjects := []sentry.SentryProject{
			{ID: "1", Slug: "project-1", Name: "Project One"},
		}
		paginatedProjects := []sentry.SentryProject{
			{ID: "1", Slug: "project-1", Name: "Project One"},
			{ID: "2", Slug: "project-2", Name: "Project Two"},
		}

		doer := &sequencedDoer{
			responses: []*http.Response{
				createMockResponse(http.StatusOK, unpaginatedProjects, nil),
				createMockResponse(http.StatusOK, paginatedProjects, nil),
			},
		}

		client, err := sentry.NewSentryClient(testBaseURL, testOrgSlug, testAuthToken, doer)
		require.NoError(t, err)

		// 1. Fetch unpaginated (withPagination=false)
		res1, err := client.GetProjects(testOrgSlug, false, false)
		require.NoError(t, err)
		require.Len(t, res1, 1)
		assert.Len(t, doer.requests, 1)

		// 2. Fetch paginated (withPagination=true) - should miss and fetch from Sentry
		res2, err := client.GetProjects(testOrgSlug, true, false)
		require.NoError(t, err)
		require.Len(t, res2, 2)
		assert.Len(t, doer.requests, 2)

		// 3. Subsequent unpaginated call should return unpaginated cache without network request
		res3, err := client.GetProjects(testOrgSlug, false, false)
		require.NoError(t, err)
		require.Len(t, res3, 1)
		assert.Len(t, doer.requests, 2)

		// 4. Subsequent paginated call should return paginated cache without network request
		res4, err := client.GetProjects(testOrgSlug, true, false)
		require.NoError(t, err)
		require.Len(t, res4, 2)
		assert.Len(t, doer.requests, 2)
	})

	t.Run("bypassCache forces fresh network request and refreshes cache", func(t *testing.T) {
		initialProjects := []sentry.SentryProject{
			{ID: "1", Slug: "initial", Name: "Initial Project"},
		}
		refreshedProjects := []sentry.SentryProject{
			{ID: "1", Slug: "refreshed", Name: "Refreshed Project"},
		}

		doer := &sequencedDoer{
			responses: []*http.Response{
				createMockResponse(http.StatusOK, initialProjects, nil),
				createMockResponse(http.StatusOK, refreshedProjects, nil),
			},
		}

		client, err := sentry.NewSentryClient(testBaseURL, testOrgSlug, testAuthToken, doer)
		require.NoError(t, err)

		// First call (withPagination=true)
		res1, err := client.GetProjects(testOrgSlug, true, false)
		require.NoError(t, err)
		assert.Equal(t, "initial", res1[0].Slug)
		assert.Len(t, doer.requests, 1)

		// Second call with bypassCache=true forces network fetch
		res2, err := client.GetProjects(testOrgSlug, true, true)
		require.NoError(t, err)
		assert.Equal(t, "refreshed", res2[0].Slug)
		assert.Len(t, doer.requests, 2)

		// Third call without bypassCache returns refreshed cached value
		res3, err := client.GetProjects(testOrgSlug, true, false)
		require.NoError(t, err)
		assert.Equal(t, "refreshed", res3[0].Slug)
		assert.Len(t, doer.requests, 2)
	})

	t.Run("defaults to sc.OrgSlug when organizationSlug is empty", func(t *testing.T) {
		projects := []sentry.SentryProject{
			{ID: "1", Slug: "default-org-proj", Name: "Default Org Project"},
		}
		doer := &sequencedDoer{
			responses: []*http.Response{
				createMockResponse(http.StatusOK, projects, nil),
			},
		}

		client, err := sentry.NewSentryClient(testBaseURL, "default-configured-org", testAuthToken, doer)
		require.NoError(t, err)

		res, err := client.GetProjects("", false, false)
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, testBaseURL+"/api/0/organizations/default-configured-org/projects/", doer.requests[0].URL.String())
	})
}
