package service

import (
	"fmt"
	"net/url"

	"github.com/notawar/mobius/internal/server/mobius"
)

// ApplyQueries sends the list of Queries to be applied (upserted) to the
// Mobius instance.
func (c *Client) ApplyQueries(specs []*mobius.QuerySpec) error {
	req := applyQuerySpecsRequest{Specs: specs}
	verb, path := "POST", "/api/latest/mobius/spec/queries"
	var responseBody applyQuerySpecsResponse
	return c.authenticatedRequest(req, verb, path, &responseBody)
}

// GetQuerySpec returns the query spec of a query by its team+name.
func (c *Client) GetQuerySpec(teamID *uint, name string) (*mobius.QuerySpec, error) {
	verb, path := "GET", "/api/latest/mobius/spec/queries/"+url.PathEscape(name)
	query := url.Values{}
	if teamID != nil {
		query.Set("team_id", fmt.Sprint(*teamID))
	}
	var responseBody getQuerySpecResponse
	err := c.authenticatedRequestWithQuery(nil, verb, path, &responseBody, query.Encode())
	return responseBody.Spec, err
}

// GetQueries retrieves the list of all Queries.
func (c *Client) GetQueries(teamID *uint, name *string) ([]mobius.Query, error) {
	verb, path := "GET", "/api/latest/mobius/queries"
	query := url.Values{}
	if teamID != nil {
		query.Set("team_id", fmt.Sprint(*teamID))
	}
	if name != nil {
		query.Set("query", *name)
	}
	var responseBody listQueriesResponse
	err := c.authenticatedRequestWithQuery(nil, verb, path, &responseBody, query.Encode())
	if err != nil {
		return nil, err
	}
	return responseBody.Queries, nil
}

// DeleteQuery deletes the query with the matching name.
func (c *Client) DeleteQuery(name string) error {
	verb, path := "DELETE", "/api/latest/mobius/queries/"+url.PathEscape(name)
	var responseBody deleteQueryResponse
	return c.authenticatedRequest(nil, verb, path, &responseBody)
}

// DeleteQueries deletes several queries.
func (c *Client) DeleteQueries(ids []uint) error {
	req := deleteQueriesRequest{IDs: ids}
	verb, path := "POST", "/api/latest/mobius/queries/delete"
	var responseBody deleteQueriesResponse
	return c.authenticatedRequest(req, verb, path, &responseBody)
}
