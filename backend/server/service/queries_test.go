package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/notawar/mobius/backend/server/contexts/viewer"
	"github.com/notawar/mobius/backend/server/datastore/mysql"
	"github.com/notawar/mobius/backend/server/mobius"
	"github.com/notawar/mobius/backend/server/mock"
	"github.com/notawar/mobius/backend/server/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryPayloadValidationCreate(t *testing.T) {
	ds := new(mock.Store)
	ds.NewQueryFunc = func(ctx context.Context, query *mobius.Query, opts ...mobius.OptionalArg) (*mobius.Query, error) {
		return query, nil
	}
	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{}, nil
	}
	ds.NewActivityFunc = func(
		ctx context.Context, user *mobius.User, activity mobius.ActivityDetails, details []byte, createdAt time.Time,
	) error {
		act, ok := activity.(mobius.ActivityTypeCreatedSavedQuery)
		assert.True(t, ok)
		assert.NotEmpty(t, act.Name)
		return nil
	}
	svc, ctx := newTestService(t, ds, nil, nil)

	testCases := []struct {
		name         string
		queryPayload mobius.QueryPayload
		shouldErr    bool
	}{
		{
			"All valid",
			mobius.QueryPayload{
				Name:     ptr.String("test query"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("snapshot"),
				Platform: ptr.String(""),
			},
			false,
		},
		{
			"Invalid  - empty string name",
			mobius.QueryPayload{
				Name:     ptr.String(""),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("snapshot"),
				Platform: ptr.String(""),
			},
			true,
		},
		{
			"Empty SQL",
			mobius.QueryPayload{
				Name:     ptr.String("bad sql"),
				Query:    ptr.String(""),
				Logging:  ptr.String("snapshot"),
				Platform: ptr.String(""),
			},
			true,
		},
		{
			"Invalid logging",
			mobius.QueryPayload{
				Name:     ptr.String("bad logging"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("hopscotch"),
				Platform: ptr.String(""),
			},
			true,
		},
		{
			"Unsupported platform",
			mobius.QueryPayload{
				Name:     ptr.String("invalid platform"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("differential"),
				Platform: ptr.String("charles"),
			},
			true,
		},
		{
			"Missing comma",
			mobius.QueryPayload{
				Name:     ptr.String("invalid platform"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("differential"),
				Platform: ptr.String("darwin windows"),
			},
			true,
		},
		{
			"Unsupported platform 'sphinx' ",
			mobius.QueryPayload{
				Name:     ptr.String("invalid platform"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("differential"),
				Platform: ptr.String("darwin,windows,sphinx"),
			},
			true,
		},
	}

	testAdmin := mobius.User{
		ID:         1,
		Teams:      []mobius.UserTeam{},
		GlobalRole: ptr.String(mobius.RoleAdmin),
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			viewerCtx := viewer.NewContext(ctx, viewer.Viewer{User: &testAdmin})
			query, err := svc.NewQuery(viewerCtx, tt.queryPayload)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Nil(t, query)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, query)
			}
		})
	}
}

// similar for modify
func TestQueryPayloadValidationModify(t *testing.T) {
	ds := new(mock.Store)
	ds.QueryFunc = func(ctx context.Context, id uint) (*mobius.Query, error) {
		return &mobius.Query{
			ID:             id,
			Name:           "mock saved query",
			Description:    "some desc",
			Query:          "select 1;",
			Platform:       "",
			Saved:          true,
			ObserverCanRun: false,
		}, nil
	}
	ds.SaveQueryFunc = func(ctx context.Context, query *mobius.Query, shouldDiscardResults bool, shouldDeleteStats bool) error {
		assert.NotEmpty(t, query)
		return nil
	}

	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{}, nil
	}
	ds.NewActivityFunc = func(
		ctx context.Context, user *mobius.User, activity mobius.ActivityDetails, details []byte, createdAt time.Time,
	) error {
		act, ok := activity.(mobius.ActivityTypeEditedSavedQuery)
		assert.True(t, ok)
		assert.NotEmpty(t, act.Name)
		return nil
	}

	svc, ctx := newTestService(t, ds, nil, nil)

	testCases := []struct {
		name         string
		queryPayload mobius.QueryPayload
		shouldErr    bool
	}{
		{
			"All valid",
			mobius.QueryPayload{
				Name:     ptr.String("updated test query"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("snapshot"),
				Platform: ptr.String(""),
			},
			false,
		},
		{
			"Invalid  - empty string name",
			mobius.QueryPayload{
				Name:     ptr.String(""),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("snapshot"),
				Platform: ptr.String(""),
			},
			true,
		},
		{
			"Empty SQL",
			mobius.QueryPayload{
				Name:     ptr.String("bad sql"),
				Query:    ptr.String(""),
				Logging:  ptr.String("snapshot"),
				Platform: ptr.String(""),
			},
			true,
		},
		{
			"Invalid logging",
			mobius.QueryPayload{
				Name:     ptr.String("bad logging"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("hopscotch"),
				Platform: ptr.String(""),
			},
			true,
		},
		{
			"Unsupported platform",
			mobius.QueryPayload{
				Name:     ptr.String("invalid platform"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("differential"),
				Platform: ptr.String("charles"),
			},
			true,
		},
		{
			"Missing comma delimeter in platform string",
			mobius.QueryPayload{
				Name:     ptr.String("invalid platform"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("differential"),
				Platform: ptr.String("darwin windows"),
			},
			true,
		},
		{
			"Unsupported platform 2",
			mobius.QueryPayload{
				Name:     ptr.String("invalid platform"),
				Query:    ptr.String("select 1"),
				Logging:  ptr.String("differential"),
				Platform: ptr.String("darwin,windows,sphinx"),
			},
			true,
		},
	}

	testAdmin := mobius.User{
		ID:         1,
		Teams:      []mobius.UserTeam{},
		GlobalRole: ptr.String(mobius.RoleAdmin),
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			viewerCtx := viewer.NewContext(ctx, viewer.Viewer{User: &testAdmin})
			_, err := svc.ModifyQuery(viewerCtx, 1, tt.queryPayload)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQueryAuth(t *testing.T) {
	ds := new(mock.Store)
	svc, ctx := newTestService(t, ds, nil, nil)

	team := mobius.Team{
		ID:   1,
		Name: "Foobar",
	}

	team2 := mobius.Team{
		ID:   2,
		Name: "Barfoo",
	}

	teamAdmin := &mobius.User{
		ID: 42,
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team.ID},
				Role: mobius.RoleAdmin,
			},
		},
	}
	teamMaintainer := &mobius.User{
		ID: 43,
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team.ID},
				Role: mobius.RoleMaintainer,
			},
		},
	}
	teamObserver := &mobius.User{
		ID: 44,
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team.ID},
				Role: mobius.RoleObserver,
			},
		},
	}
	teamObserverPlus := &mobius.User{
		ID: 45,
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team.ID},
				Role: mobius.RoleObserverPlus,
			},
		},
	}
	teamGitOps := &mobius.User{
		ID: 46,
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team.ID},
				Role: mobius.RoleGitOps,
			},
		},
	}
	globalQuery := mobius.Query{
		ID:     99,
		Name:   "global query",
		TeamID: nil,
	}
	teamQuery := mobius.Query{
		ID:     88,
		Name:   "team query",
		TeamID: ptr.Uint(team.ID),
	}
	team2Query := mobius.Query{
		ID:     77,
		Name:   "team2 query",
		TeamID: ptr.Uint(team2.ID),
	}
	queriesMap := map[uint]mobius.Query{
		globalQuery.ID: globalQuery,
		teamQuery.ID:   teamQuery,
		team2Query.ID:  team2Query,
	}

	ds.TeamFunc = func(ctx context.Context, tid uint) (*mobius.Team, error) {
		if tid == team.ID {
			return &team, nil
		} else if tid == team2.ID {
			return &team2, nil
		}
		return nil, newNotFoundError()
	}

	ds.TeamByNameFunc = func(ctx context.Context, name string) (*mobius.Team, error) {
		if name == team.Name {
			return &team, nil
		} else if name == team2.Name {
			return &team2, nil
		}
		return nil, newNotFoundError()
	}
	ds.NewQueryFunc = func(ctx context.Context, query *mobius.Query, opts ...mobius.OptionalArg) (*mobius.Query, error) {
		return query, nil
	}
	ds.QueryByNameFunc = func(ctx context.Context, teamID *uint, name string) (*mobius.Query, error) {
		if teamID == nil && name == "global query" { //nolint:gocritic // ignore ifElseChain
			return &globalQuery, nil
		} else if teamID != nil && *teamID == team.ID && name == "team query" {
			return &teamQuery, nil
		} else if teamID != nil && *teamID == team2.ID && name == "team2 query" {
			return &team2Query, nil
		}
		return nil, newNotFoundError()
	}
	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{}, nil
	}
	ds.NewActivityFunc = func(
		ctx context.Context, user *mobius.User, activity mobius.ActivityDetails, details []byte, createdAt time.Time,
	) error {
		return nil
	}
	ds.QueryFunc = func(ctx context.Context, id uint) (*mobius.Query, error) {
		if id == 99 { //nolint:gocritic // ignore ifElseChain
			return &globalQuery, nil
		} else if id == 88 {
			return &teamQuery, nil
		} else if id == 77 {
			return &team2Query, nil
		}
		return nil, newNotFoundError()
	}

	ds.ResultCountForQueryFunc = func(ctx context.Context, queryID uint) (int, error) {
		return 0, nil
	}

	ds.SaveQueryFunc = func(ctx context.Context, query *mobius.Query, shouldDiscardResults bool, shouldDeleteStats bool) error {
		return nil
	}
	ds.DeleteQueryFunc = func(ctx context.Context, teamID *uint, name string) error {
		return nil
	}
	ds.DeleteQueriesFunc = func(ctx context.Context, ids []uint) (uint, error) {
		return 0, nil
	}
	ds.ListQueriesFunc = func(ctx context.Context, opts mobius.ListQueryOptions) ([]*mobius.Query, int, *mobius.PaginationMetadata, error) {
		return nil, 0, nil, nil
	}
	ds.ApplyQueriesFunc = func(ctx context.Context, authID uint, queries []*mobius.Query, queriesToDiscardResults map[uint]struct{}) error {
		return nil
	}

	testCases := []struct {
		name            string
		user            *mobius.User
		qid             uint
		shouldFailWrite bool
		shouldFailRead  bool
		shouldFailNew   bool
	}{
		{
			"global admin and global query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleAdmin)},
			globalQuery.ID,
			false,
			false,
			false,
		},
		{
			"global admin and team query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleAdmin)},
			teamQuery.ID,
			false,
			false,
			false,
		},
		{
			"global maintainer and global query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleMaintainer)},
			globalQuery.ID,
			false,
			false,
			false,
		},
		{
			"global maintainer and team query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleMaintainer)},
			teamQuery.ID,
			false,
			false,
			false,
		},
		{
			"global observer and global query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleObserver)},
			globalQuery.ID,
			true,
			false,
			true,
		},
		{
			"global observer and team query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleObserver)},
			teamQuery.ID,
			true,
			false,
			true,
		},
		{
			"global observer+ and global query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleObserverPlus)},
			globalQuery.ID,
			true,
			false,
			true,
		},
		{
			"global observer+ and team query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleObserverPlus)},
			teamQuery.ID,
			true,
			false,
			true,
		},
		{
			"global gitops and global query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleGitOps)},
			globalQuery.ID,
			false,
			false,
			false,
		},
		{
			"global gitops and team query",
			&mobius.User{GlobalRole: ptr.String(mobius.RoleGitOps)},
			teamQuery.ID,
			false,
			false,
			false,
		},
		{
			"team admin and global query",
			teamAdmin,
			globalQuery.ID,
			true,
			false,
			true,
		},
		{
			"team admin and team query",
			teamAdmin,
			teamQuery.ID,
			false,
			false,
			false,
		},
		{
			"team admin and team2 query",
			teamAdmin,
			team2Query.ID,
			true,
			true,
			true,
		},
		{
			"team maintainer and global query",
			teamMaintainer,
			globalQuery.ID,
			true,
			false,
			true,
		},
		{
			"team maintainer and team query",
			teamMaintainer,
			teamQuery.ID,
			false,
			false,
			false,
		},
		{
			"team maintainer and team2 query",
			teamMaintainer,
			team2Query.ID,
			true,
			true,
			true,
		},
		{
			"team observer and global query",
			teamObserver,
			globalQuery.ID,
			true,
			false,
			true,
		},
		{
			"team observer and team query",
			teamObserver,
			teamQuery.ID,
			true,
			false,
			true,
		},
		{
			"team observer and team2 query",
			teamObserver,
			team2Query.ID,
			true,
			true,
			true,
		},
		{
			"team observer+ and global query",
			teamObserverPlus,
			globalQuery.ID,
			true,
			false,
			true,
		},
		{
			"team observer+ and team query",
			teamObserverPlus,
			teamQuery.ID,
			true,
			false,
			true,
		},
		{
			"team observer+ and team2 query",
			teamObserverPlus,
			team2Query.ID,
			true,
			true,
			true,
		},
		{
			"team gitops and global query",
			teamGitOps,
			globalQuery.ID,
			true,
			true,
			true,
		},
		{
			"team gitops and team query",
			teamGitOps,
			teamQuery.ID,
			false,
			false,
			false,
		},
		{
			"team gitops and team2 query",
			teamGitOps,
			team2Query.ID,
			true,
			true,
			true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := viewer.NewContext(ctx, viewer.Viewer{User: tt.user})

			query := queriesMap[tt.qid]

			_, err := svc.NewQuery(ctx, mobius.QueryPayload{
				Name:   ptr.String("name"),
				Query:  ptr.String("select 1"),
				TeamID: query.TeamID,
			})
			checkAuthErr(t, tt.shouldFailNew, err)

			_, err = svc.ModifyQuery(ctx, tt.qid, mobius.QueryPayload{})
			checkAuthErr(t, tt.shouldFailWrite, err)

			err = svc.DeleteQuery(ctx, query.TeamID, query.Name)
			checkAuthErr(t, tt.shouldFailWrite, err)

			err = svc.DeleteQueryByID(ctx, tt.qid)
			checkAuthErr(t, tt.shouldFailWrite, err)

			_, err = svc.DeleteQueries(ctx, []uint{tt.qid})
			checkAuthErr(t, tt.shouldFailWrite, err)

			_, err = svc.GetQuery(ctx, tt.qid)
			checkAuthErr(t, tt.shouldFailRead, err)

			_, err = svc.QueryReportIsClipped(ctx, tt.qid, mobius.DefaultMaxQueryReportRows)
			checkAuthErr(t, tt.shouldFailRead, err)

			_, _, _, err = svc.ListQueries(ctx, mobius.ListOptions{}, query.TeamID, nil, false, nil)
			checkAuthErr(t, tt.shouldFailRead, err)

			teamName := ""
			if query.TeamID != nil && *query.TeamID == team.ID {
				teamName = team.Name
			} else if query.TeamID != nil && *query.TeamID == team2.ID {
				teamName = team2.Name
			}

			err = svc.ApplyQuerySpecs(ctx, []*mobius.QuerySpec{{
				Name:     query.Name,
				Query:    "SELECT 1",
				TeamName: teamName,
			}})
			checkAuthErr(t, tt.shouldFailWrite, err)

			_, err = svc.GetQuerySpecs(ctx, query.TeamID)
			checkAuthErr(t, tt.shouldFailRead, err)

			_, err = svc.GetQuerySpec(ctx, query.TeamID, query.Name)
			checkAuthErr(t, tt.shouldFailRead, err)
		})
	}
}

func TestQueryReportIsClipped(t *testing.T) {
	ds := new(mock.Store)
	svc, ctx := newTestService(t, ds, nil, nil)
	viewerCtx := viewer.NewContext(ctx, viewer.Viewer{User: &mobius.User{
		ID:         1,
		GlobalRole: ptr.String(mobius.RoleAdmin),
	}})

	ds.QueryFunc = func(ctx context.Context, queryID uint) (*mobius.Query, error) {
		return &mobius.Query{}, nil
	}
	ds.ResultCountForQueryFunc = func(ctx context.Context, queryID uint) (int, error) {
		return 0, nil
	}

	isClipped, err := svc.QueryReportIsClipped(viewerCtx, 1, mobius.DefaultMaxQueryReportRows)
	require.NoError(t, err)
	require.False(t, isClipped)

	ds.ResultCountForQueryFunc = func(ctx context.Context, queryID uint) (int, error) {
		return mobius.DefaultMaxQueryReportRows, nil
	}

	isClipped, err = svc.QueryReportIsClipped(viewerCtx, 1, mobius.DefaultMaxQueryReportRows)
	require.NoError(t, err)
	require.True(t, isClipped)
}

func TestQueryReportReturnsNilIfDiscardDataIsTrue(t *testing.T) {
	ds := new(mock.Store)
	svc, ctx := newTestService(t, ds, nil, nil)
	viewerCtx := viewer.NewContext(ctx, viewer.Viewer{User: &mobius.User{
		ID:         1,
		GlobalRole: ptr.String(mobius.RoleAdmin),
	}})

	ds.QueryFunc = func(ctx context.Context, queryID uint) (*mobius.Query, error) {
		return &mobius.Query{
			DiscardData: true,
		}, nil
	}
	ds.QueryResultRowsFunc = func(ctx context.Context, queryID uint, opts mobius.TeamFilter) ([]*mobius.ScheduledQueryResultRow, error) {
		return []*mobius.ScheduledQueryResultRow{
			{
				QueryID:     1,
				HostID:      1,
				Data:        ptr.RawMessage(json.RawMessage(`{"foo": "bar"}`)),
				LastFetched: time.Now(),
			},
		}, nil
	}

	results, reportClipped, err := svc.GetQueryReportResults(viewerCtx, 1, nil)
	require.NoError(t, err)
	require.Nil(t, results)
	require.False(t, reportClipped)
}

func TestInheritedQueryReportTeamPermissions(t *testing.T) {
	ds := mysql.CreateMySQLDS(t)
	defer ds.Close()

	svc, ctx := newTestService(t, ds, nil, nil)

	team1, err := ds.NewTeam(ctx, &mobius.Team{
		ID:          42,
		Name:        "team1",
		Description: "desc team1",
	})
	require.NoError(t, err)
	team2, err := ds.NewTeam(ctx, &mobius.Team{
		Name:        "team2",
		Description: "desc team2",
	})
	require.NoError(t, err)

	hostTeam2, err := ds.NewHost(ctx, &mobius.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		PolicyUpdatedAt: time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         ptr.String("1"),
		UUID:            "1",
		ComputerName:    "Foo Local",
		Hostname:        "foo.local",
		OsqueryHostID:   ptr.String("1"),
		PrimaryIP:       "192.168.1.1",
		PrimaryMac:      "30-65-EC-6F-C4-61",
		Platform:        "darwin",
	})
	require.NoError(t, err)
	err = ds.AddHostsToTeam(ctx, &team2.ID, []uint{hostTeam2.ID})
	require.NoError(t, err)

	hostTeam1, err := ds.NewHost(ctx, &mobius.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		PolicyUpdatedAt: time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         ptr.String("42"),
		UUID:            "42",
		ComputerName:    "bar Local",
		Hostname:        "bar.local",
		OsqueryHostID:   ptr.String("42"),
		PrimaryIP:       "192.168.1.2",
		PrimaryMac:      "30-65-EC-6F-C4-62",
		Platform:        "darwin",
	})
	require.NoError(t, err)
	err = ds.AddHostsToTeam(ctx, &team1.ID, []uint{hostTeam1.ID})
	require.NoError(t, err)

	globalQuery, err := ds.NewQuery(ctx, &mobius.Query{
		ID:      77,
		Name:    "team2 query",
		TeamID:  nil,
		Query:   "select * from usb_devices;",
		Logging: mobius.LoggingSnapshot,
	})
	require.NoError(t, err)
	// Insert initial Result Rows
	mockTime := time.Now().UTC().Truncate(time.Second)
	host2Row := []*mobius.ScheduledQueryResultRow{
		{
			QueryID:     globalQuery.ID,
			HostID:      hostTeam2.ID,
			LastFetched: mockTime,
			Data:        ptr.RawMessage([]byte(`{"model": "USB Keyboard", "vendor": "Apple Inc."}`)),
		},
	}
	err = ds.OverwriteQueryResultRows(ctx, host2Row, mobius.DefaultMaxQueryReportRows)
	require.NoError(t, err)
	host1Row := []*mobius.ScheduledQueryResultRow{
		{
			QueryID:     globalQuery.ID,
			HostID:      hostTeam1.ID,
			LastFetched: mockTime,
			Data:        ptr.RawMessage([]byte(`{"model": "USB Mouse", "vendor": "Apple Inc."}`)),
		},
	}
	err = ds.OverwriteQueryResultRows(ctx, host1Row, mobius.DefaultMaxQueryReportRows)
	require.NoError(t, err)

	team2Admin := &mobius.User{
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team2.ID},
				Role: mobius.RoleAdmin,
			},
		},
	}

	queryReportResults, _, err := svc.GetQueryReportResults(viewer.NewContext(ctx, viewer.Viewer{User: team2Admin}), globalQuery.ID, &team2.ID)
	require.NoError(t, err)
	require.Len(t, queryReportResults, 1)

	// team admins requesting query results filtered to not-their-team should get no rows back

	teamAdmin := &mobius.User{
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team1.ID},
				Role: mobius.RoleAdmin,
			},
		},
	}
	teamMaintainer := &mobius.User{
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team1.ID},
				Role: mobius.RoleMaintainer,
			},
		},
	}
	teamObserver := &mobius.User{
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team1.ID},
				Role: mobius.RoleObserver,
			},
		},
	}
	teamObserverPlus := &mobius.User{
		Teams: []mobius.UserTeam{
			{
				Team: mobius.Team{ID: team1.ID},
				Role: mobius.RoleObserverPlus,
			},
		},
	}

	testCases := []struct {
		name string
		user *mobius.User
	}{
		{
			name: "team admin",
			user: teamAdmin,
		},
		{
			name: "team maintainer",
			user: teamMaintainer,
		},
		{
			name: "team observer",
			user: teamObserver,
		},
		{
			name: "team observer+",
			user: teamObserverPlus,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			queryReportResults, _, err := svc.GetQueryReportResults(viewer.NewContext(ctx, viewer.Viewer{User: tt.user}), globalQuery.ID, &team2.ID)
			require.NoError(t, err)
			require.Len(t, queryReportResults, 0)
		})
	}
}

func TestComparePlatforms(t *testing.T) {
	for _, tc := range []struct {
		name     string
		p1       string
		p2       string
		expected bool
	}{
		{
			name:     "equal single value",
			p1:       "linux",
			p2:       "linux",
			expected: true,
		},
		{
			name:     "different single value",
			p1:       "macos",
			p2:       "linux",
			expected: false,
		},
		{
			name:     "equal multiple values",
			p1:       "linux,windows",
			p2:       "linux,windows",
			expected: true,
		},
		{
			name:     "equal multiple values out of order",
			p1:       "linux,windows",
			p2:       "windows,linux",
			expected: true,
		},
		{
			name:     "different multiple values",
			p1:       "linux,windows",
			p2:       "linux,windows,darwin",
			expected: false,
		},
		{
			name:     "no values set",
			p1:       "",
			p2:       "",
			expected: true,
		},
		{
			name:     "no values set",
			p1:       "",
			p2:       "linux",
			expected: false,
		},
		{
			name:     "single and multiple values",
			p1:       "linux",
			p2:       "windows,linux",
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual := comparePlatforms(tc.p1, tc.p2)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestApplyQuerySpec(t *testing.T) {
	ds := new(mock.Store)
	ds.NewQueryFunc = func(ctx context.Context, query *mobius.Query, opts ...mobius.OptionalArg) (*mobius.Query, error) {
		return query, nil
	}
	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{}, nil
	}
	ds.NewActivityFunc = func(
		ctx context.Context, user *mobius.User, activity mobius.ActivityDetails, details []byte, createdAt time.Time,
	) error {
		return nil
	}
	ds.QueryByNameFunc = func(ctx context.Context, teamID *uint, name string) (*mobius.Query, error) {
		return nil, newNotFoundError()
	}
	ds.ApplyQueriesFunc = func(ctx context.Context, authID uint, queries []*mobius.Query, queriesToDiscardResults map[uint]struct{}) error {
		return nil
	}
	ds.LabelsByNameFunc = func(ctx context.Context, names []string) (map[string]*mobius.Label, error) {
		labels := make(map[string]*mobius.Label, len(names))
		for _, name := range names {
			if name == "foo" {
				labels["foo"] = &mobius.Label{
					Name: "foo",
					ID:   1,
				}
			}
		}
		return labels, nil
	}

	svc, ctx := newTestService(t, ds, nil, nil)

	testAdmin := mobius.User{
		ID:         1,
		Teams:      []mobius.UserTeam{},
		GlobalRole: ptr.String(mobius.RoleAdmin),
	}
	viewerCtx := viewer.NewContext(ctx, viewer.Viewer{User: &testAdmin})

	// Test that a query spec with a label that exists doesn't return an error
	err := svc.ApplyQuerySpecs(viewerCtx, []*mobius.QuerySpec{
		{
			Name:             "test query",
			Query:            "select 1",
			LabelsIncludeAny: []string{"foo"},
			Platform:         "darwin,windows",
		},
	})
	require.NoError(t, err)

	// Test that a query spec with a label that doesn't exist returns an error.
	err = svc.ApplyQuerySpecs(viewerCtx, []*mobius.QuerySpec{
		{
			Name:             "test query",
			Query:            "select 1",
			LabelsIncludeAny: []string{"foo", "bar"},
			Platform:         "darwin,windows",
		},
	})
	assert.Error(t, err)
}
