package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/notawar/mobius/server/calendar"
	"github.com/notawar/mobius/server/config"
	"github.com/notawar/mobius/server/datastore/redis/redistest"
	"github.com/notawar/mobius/server/mobius"
	"github.com/notawar/mobius/server/mock"
	"github.com/notawar/mobius/server/ptr"
	"github.com/notawar/mobius/server/service/redis_lock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultCalendarConfig = config.CalendarConfig{Periodicity: 5 * time.Minute}

func TestGetPreferredCalendarEventDate(t *testing.T) {
	t.Parallel()
	date := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
	for _, tc := range []struct {
		name      string
		year      int
		month     time.Month
		daysStart int
		daysEnd   int

		expected time.Time
	}{
		{
			name:      "March 2024 (before 1st Tuesday)",
			year:      2024,
			month:     3,
			daysStart: 1,
			daysEnd:   5,

			expected: date(2024, 3, 5),
		},
		{
			name:      "March 2024 (past 1st Tuesday)",
			year:      2024,
			month:     3,
			daysStart: 6,
			daysEnd:   12,

			expected: date(2024, 3, 12),
		},
		{
			name:      "April 2024 (before 3rd Tuesday)",
			year:      2024,
			month:     4,
			daysStart: 10,
			daysEnd:   16,

			expected: date(2024, 4, 16),
		},
		{
			name:      "April 2024 (after 3rd Tuesday)",
			year:      2024,
			month:     4,
			daysStart: 17,
			daysEnd:   23,

			expected: date(2024, 4, 23),
		},
		{
			name:      "May 2024 (before last Tuesday)",
			year:      2024,
			month:     5,
			daysStart: 22,
			daysEnd:   28,

			expected: date(2024, 5, 28),
		},
		{
			name:      "May 2024 (after last Tuesday)",
			year:      2024,
			month:     5,
			daysStart: 29,
			daysEnd:   31,

			expected: date(2024, 6, 4),
		},
		{
			name:      "Dec 2025 (before last Tuesday)",
			year:      2025,
			month:     12,
			daysStart: 24,
			daysEnd:   30,

			expected: date(2025, 12, 30),
		},
		{
			name:      "Dec 2025 (after last Tuesday)",
			year:      2025,
			month:     12,
			daysStart: 31,
			daysEnd:   31,

			expected: date(2026, 1, 6),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for day := tc.daysStart; day <= tc.daysEnd; day++ {
				actual := getPreferredCalendarEventDate(tc.year, tc.month, day)
				require.NotEqual(t, actual.Weekday(), time.Saturday)
				require.NotEqual(t, actual.Weekday(), time.Sunday)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

// TestEventForDifferentHost tests case when event exists, but for a different host. Nothing should happen.
// The old event will eventually be cleaned up by the cleanup job, and afterward a new event will be created.
func TestEventForDifferentHost(t *testing.T) {
	t.Parallel()
	ds := new(mock.Store)
	ctx := context.Background()
	logger := kitlog.With(kitlog.NewLogfmtLogger(os.Stdout))
	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{
			Integrations: mobius.Integrations{
				GoogleCalendar: []*mobius.GoogleCalendarIntegration{
					{},
				},
			},
		}, nil
	}
	teamID1 := uint(1)
	ds.ListTeamsFunc = func(ctx context.Context, filter mobius.TeamFilter, opt mobius.ListOptions) ([]*mobius.Team, error) {
		return []*mobius.Team{
			{
				ID: teamID1,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable: true,
						},
					},
				},
			},
		}, nil
	}
	policyID1 := uint(10)
	ds.GetCalendarPoliciesFunc = func(ctx context.Context, teamID uint) ([]mobius.PolicyCalendarData, error) {
		require.Equal(t, teamID1, teamID)
		return []mobius.PolicyCalendarData{
			{
				ID:   policyID1,
				Name: "Policy 1",
			},
		}, nil
	}
	hostID1 := uint(100)
	hostID2 := uint(101)
	userEmail1 := "user@example.com"
	ds.GetTeamHostsPolicyMembershipsFunc = func(
		ctx context.Context, domain string, teamID uint, policyIDs []uint, _ *uint,
	) ([]mobius.HostPolicyMembershipData, error) {
		require.Equal(t, teamID1, teamID)
		require.Equal(t, []uint{policyID1}, policyIDs)
		return []mobius.HostPolicyMembershipData{
			{
				HostID:           hostID1,
				Email:            userEmail1,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d", policyID1),
			},
		}, nil
	}
	// Return an existing event, but for a different host
	eventTime := time.Now().Add(time.Hour)
	ds.GetHostCalendarEventByEmailFunc = func(ctx context.Context, email string) (*mobius.HostCalendarEvent, *mobius.CalendarEvent, error) {
		require.Equal(t, userEmail1, email)
		calEvent := &mobius.CalendarEvent{
			ID:        1,
			Email:     email,
			StartTime: eventTime,
			EndTime:   eventTime,
		}
		hcEvent := &mobius.HostCalendarEvent{
			ID:              1,
			HostID:          hostID2,
			CalendarEventID: 1,
			WebhookStatus:   mobius.CalendarWebhookStatusNone,
		}
		return hcEvent, calEvent, nil
	}

	pool := redistest.SetupRedis(t, t.Name(), false, false, false)
	err := cronCalendarEvents(ctx, ds, redis_lock.NewLock(pool), defaultCalendarConfig, logger)
	require.NoError(t, err)
}

func TestCalendarEventsMultipleHosts(t *testing.T) {
	ds := new(mock.Store)
	ctx := context.Background()
	logger := kitlog.With(kitlog.NewLogfmtLogger(os.Stdout))
	t.Cleanup(func() {
		calendar.ClearMockEvents()
		calendar.ClearMockChannels()
	})

	//
	// Test setup
	//
	//	team1:
	//
	//	policyID1 (calendar)
	//	policyID2 (calendar)
	//
	// 	hostID1 has user1@example.com not passing policies.
	//	hostID2 has user2@example.com passing policies.
	//	hostID3 does not have example.com email and is not passing policies.
	//	hostID4 does not have example.com email and is passing policies.
	//

	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{
			Integrations: mobius.Integrations{
				GoogleCalendar: []*mobius.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: map[string]string{
							mobius.GoogleCalendarEmail: "calendar-mock@example.com",
						},
					},
				},
			},
		}, nil
	}

	teamID1 := uint(1)
	ds.ListTeamsFunc = func(ctx context.Context, filter mobius.TeamFilter, opt mobius.ListOptions) ([]*mobius.Team, error) {
		return []*mobius.Team{
			{
				ID: teamID1,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable:     true,
							WebhookURL: "https://foo.example.com",
						},
					},
				},
			},
		}, nil
	}

	policyID1 := uint(10)
	policyID2 := uint(11)
	ds.GetCalendarPoliciesFunc = func(ctx context.Context, teamID uint) ([]mobius.PolicyCalendarData, error) {
		require.Equal(t, teamID1, teamID)
		return []mobius.PolicyCalendarData{
			{
				ID:   policyID1,
				Name: "Policy 1",
			},
			{
				ID:   policyID2,
				Name: "Policy 2",
			},
		}, nil
	}

	hostID1, userEmail1 := uint(100), "user1@example.com"
	hostID2, userEmail2 := uint(101), "user2@example.com"
	hostID3 := uint(102)
	hostID4 := uint(103)

	ds.GetTeamHostsPolicyMembershipsFunc = func(
		ctx context.Context, domain string, teamID uint, policyIDs []uint, _ *uint,
	) ([]mobius.HostPolicyMembershipData, error) {
		require.Equal(t, "example.com", domain)
		require.Equal(t, teamID1, teamID)
		require.Equal(t, []uint{policyID1, policyID2}, policyIDs)
		return []mobius.HostPolicyMembershipData{
			{
				HostID:           hostID1,
				Email:            userEmail1,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d,%d", policyID1, policyID2),
			},
			{
				HostID:  hostID2,
				Email:   userEmail2,
				Passing: true,
			},
			{
				HostID:           hostID3,
				Email:            "", // because it does not belong to example.com
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d,%d", policyID1, policyID2),
			},
			{
				HostID:  hostID4,
				Email:   "", // because it does not belong to example.com
				Passing: true,
			},
		}, nil
	}
	ds.PolicyLiteFunc = func(ctx context.Context, policyID uint) (*mobius.PolicyLite, error) {
		switch policyID {
		case policyID1:
			return &mobius.PolicyLite{
				ID:          policyID1,
				Description: "Policy 1",
			}, nil
		case policyID2:
			return &mobius.PolicyLite{
				ID:          policyID2,
				Description: "Policy 2",
			}, nil
		default:
			t.Errorf("unexpected policy ID: %d", policyID)
			return nil, nil
		}
	}

	ds.GetHostCalendarEventByEmailFunc = func(ctx context.Context, email string) (*mobius.HostCalendarEvent, *mobius.CalendarEvent, error) {
		return nil, nil, notFoundErr{}
	}

	var eventsMu sync.Mutex
	calendarEvents := make(map[string]*mobius.CalendarEvent)
	hostCalendarEvents := make(map[uint]*mobius.HostCalendarEvent)

	ds.CreateOrUpdateCalendarEventFunc = func(ctx context.Context,
		uuid string,
		email string,
		startTime, endTime time.Time,
		data []byte,
		timeZone *string,
		hostID uint,
		webhookStatus mobius.CalendarWebhookStatus,
	) (*mobius.CalendarEvent, error) {
		assert.NotEmpty(t, uuid)
		require.Equal(t, hostID1, hostID)
		require.Equal(t, userEmail1, email)
		require.Equal(t, mobius.CalendarWebhookStatusNone, webhookStatus)
		require.NotEmpty(t, data)
		require.NotZero(t, startTime)
		require.NotZero(t, endTime)

		eventsMu.Lock()
		calendarEventID := uint(len(calendarEvents) + 1) //nolint:gosec // dismiss G115
		calendarEvents[email] = &mobius.CalendarEvent{
			ID:        calendarEventID,
			Email:     email,
			StartTime: startTime,
			EndTime:   endTime,
			Data:      data,
		}
		hostCalendarEventID := uint(len(hostCalendarEvents) + 1) //nolint:gosec // dismiss G115
		hostCalendarEvents[hostID] = &mobius.HostCalendarEvent{
			ID:              hostCalendarEventID,
			HostID:          hostID,
			CalendarEventID: calendarEventID,
			WebhookStatus:   webhookStatus,
		}
		eventsMu.Unlock()
		return nil, nil
	}

	pool := redistest.SetupRedis(t, t.Name(), false, false, false)
	err := cronCalendarEvents(ctx, ds, redis_lock.NewLock(pool), defaultCalendarConfig, logger)
	require.NoError(t, err)

	eventsMu.Lock()
	require.Len(t, calendarEvents, 1)
	require.Len(t, hostCalendarEvents, 1)
	eventsMu.Unlock()

	createdCalendarEvents := calendar.ListGoogleMockEvents()
	require.Len(t, createdCalendarEvents, 1)
	strings.Contains(createdCalendarEvents["1"].Description, mobius.CalendarDefaultDescription)
	strings.Contains(createdCalendarEvents["1"].Description, mobius.CalendarDefaultResolution)
}

type notFoundErr struct{}

func (n notFoundErr) IsNotFound() bool {
	return true
}

func (n notFoundErr) Error() string {
	return "not found"
}

func TestCalendarEvents1KHosts(t *testing.T) {
	t.Parallel()
	ds := new(mock.Store)
	ctx := context.Background()
	var logger kitlog.Logger
	if os.Getenv("CALENDAR_TEST_LOGGING") != "" {
		logger = kitlog.With(kitlog.NewLogfmtLogger(os.Stdout))
	} else {
		logger = kitlog.NewNopLogger()
	}
	t.Cleanup(func() {
		calendar.ClearMockEvents()
		calendar.ClearMockChannels()
	})

	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{
			Integrations: mobius.Integrations{
				GoogleCalendar: []*mobius.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: map[string]string{
							mobius.GoogleCalendarEmail: "calendar-mock@example.com",
						},
					},
				},
			},
		}, nil
	}

	teamID1 := uint(1)
	teamID2 := uint(2)
	teamID3 := uint(3)
	teamID4 := uint(4)
	teamID5 := uint(5)
	ds.ListTeamsFunc = func(ctx context.Context, filter mobius.TeamFilter, opt mobius.ListOptions) ([]*mobius.Team, error) {
		return []*mobius.Team{
			{
				ID: teamID1,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable:     true,
							WebhookURL: "https://foo.example.com",
						},
					},
				},
			},
			{
				ID: teamID2,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable:     true,
							WebhookURL: "https://foo.example.com",
						},
					},
				},
			},
			{
				ID: teamID3,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable:     true,
							WebhookURL: "https://foo.example.com",
						},
					},
				},
			},
			{
				ID: teamID4,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable:     true,
							WebhookURL: "https://foo.example.com",
						},
					},
				},
			},
			{
				ID: teamID5,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable:     true,
							WebhookURL: "https://foo.example.com",
						},
					},
				},
			},
		}, nil
	}

	policyID1 := uint(10)
	policyID2 := uint(11)
	policyID3 := uint(12)
	policyID4 := uint(13)
	policyID5 := uint(14)
	policyID6 := uint(15)
	policyID7 := uint(16)
	policyID8 := uint(17)
	policyID9 := uint(18)
	policyID10 := uint(19)
	ds.GetCalendarPoliciesFunc = func(ctx context.Context, teamID uint) ([]mobius.PolicyCalendarData, error) {
		switch teamID {
		case teamID1:
			return []mobius.PolicyCalendarData{
				{
					ID:   policyID1,
					Name: "Policy 1",
				},
				{
					ID:   policyID2,
					Name: "Policy 2",
				},
			}, nil
		case teamID2:
			return []mobius.PolicyCalendarData{
				{
					ID:   policyID3,
					Name: "Policy 3",
				},
				{
					ID:   policyID4,
					Name: "Policy 4",
				},
			}, nil
		case teamID3:
			return []mobius.PolicyCalendarData{
				{
					ID:   policyID5,
					Name: "Policy 5",
				},
				{
					ID:   policyID6,
					Name: "Policy 6",
				},
			}, nil
		case teamID4:
			return []mobius.PolicyCalendarData{
				{
					ID:   policyID7,
					Name: "Policy 7",
				},
				{
					ID:   policyID8,
					Name: "Policy 8",
				},
			}, nil
		case teamID5:
			return []mobius.PolicyCalendarData{
				{
					ID:   policyID9,
					Name: "Policy 9",
				},
				{
					ID:   policyID10,
					Name: "Policy 10",
				},
			}, nil
		default:
			return nil, notFoundErr{}
		}
	}

	hosts := make([]mobius.HostPolicyMembershipData, 0, 1000)
	for i := 0; i < 1000; i++ {
		newHost := mobius.HostPolicyMembershipData{
			Email:              fmt.Sprintf("user%d@example.com", i),
			Passing:            i%2 == 0,
			HostID:             uint(i), //nolint:gosec // dismiss G115
			HostDisplayName:    fmt.Sprintf("display_name%d", i),
			HostHardwareSerial: fmt.Sprintf("serial%d", i),
		}
		if !newHost.Passing {
			switch {
			case i >= 0 && i < 200:
				newHost.FailingPolicyIDs = fmt.Sprintf("%d,%d", policyID1, policyID2)
			case i >= 200 && i < 400:
				newHost.FailingPolicyIDs = fmt.Sprintf("%d", policyID4)
			case i >= 400 && i < 600:
				newHost.FailingPolicyIDs = fmt.Sprintf("%d", policyID5)
			case i >= 600 && i < 800:
				newHost.FailingPolicyIDs = fmt.Sprintf("%d,%d", policyID7, policyID8)
			default:
				newHost.FailingPolicyIDs = fmt.Sprintf("%d,%d", policyID9, policyID10)
			}
		}
		hosts = append(hosts, newHost)
	}
	ds.PolicyLiteFunc = func(ctx context.Context, policyID uint) (*mobius.PolicyLite, error) {
		resolution := fmt.Sprintf("Resolution for policy %d", policyID)
		return &mobius.PolicyLite{
			ID:          policyID,
			Description: fmt.Sprintf("Policy %d", policyID),
			Resolution:  &resolution,
		}, nil
	}

	ds.GetTeamHostsPolicyMembershipsFunc = func(
		ctx context.Context, domain string, teamID uint, policyIDs []uint, _ *uint,
	) ([]mobius.HostPolicyMembershipData, error) {
		var start, end int
		switch teamID {
		case teamID1:
			start, end = 0, 200
		case teamID2:
			start, end = 200, 400
		case teamID3:
			start, end = 400, 600
		case teamID4:
			start, end = 600, 800
		case teamID5:
			start, end = 800, 1000
		}
		return hosts[start:end], nil
	}

	ds.GetHostCalendarEventByEmailFunc = func(ctx context.Context, email string) (*mobius.HostCalendarEvent, *mobius.CalendarEvent, error) {
		return nil, nil, notFoundErr{}
	}

	eventsCreated := 0
	var eventsCreatedMu sync.Mutex

	eventPerHost := make(map[uint]*mobius.CalendarEvent)

	ds.CreateOrUpdateCalendarEventFunc = func(ctx context.Context,
		uuid string,
		email string,
		startTime, endTime time.Time,
		data []byte,
		timeZone *string,
		hostID uint,
		webhookStatus mobius.CalendarWebhookStatus,
	) (*mobius.CalendarEvent, error) {
		assert.NotEmpty(t, uuid)
		require.Equal(t, fmt.Sprintf("user%d@example.com", hostID), email)
		eventsCreatedMu.Lock()
		eventsCreated += 1
		eventPerHost[hostID] = &mobius.CalendarEvent{
			ID:        hostID,
			Email:     email,
			StartTime: startTime,
			EndTime:   endTime,
			Data:      data,
			UpdateCreateTimestamps: mobius.UpdateCreateTimestamps{
				CreateTimestamp: mobius.CreateTimestamp{
					CreatedAt: time.Now(),
				},
				UpdateTimestamp: mobius.UpdateTimestamp{
					UpdatedAt: time.Now(),
				},
			},
		}
		eventsCreatedMu.Unlock()
		require.Equal(t, mobius.CalendarWebhookStatusNone, webhookStatus)
		require.NotEmpty(t, data)
		require.NotZero(t, startTime)
		require.NotZero(t, endTime)
		// Currently, the returned calendar event is unused.
		return nil, nil
	}

	pool := redistest.SetupRedis(t, t.Name(), false, false, false)
	distributedLock := redis_lock.NewLock(pool)
	err := cronCalendarEvents(ctx, ds, distributedLock, defaultCalendarConfig, logger)
	require.NoError(t, err)

	createdCalendarEvents := calendar.ListGoogleMockEvents()
	require.Equal(t, eventsCreated, 500)
	require.Len(t, createdCalendarEvents, 500)

	hosts = make([]mobius.HostPolicyMembershipData, 0, 1000)
	for i := 0; i < 1000; i++ {
		hosts = append(hosts, mobius.HostPolicyMembershipData{
			Email:              fmt.Sprintf("user%d@example.com", i),
			Passing:            true,
			HostID:             uint(i), //nolint:gosec // dismiss G115
			HostDisplayName:    fmt.Sprintf("display_name%d", i),
			HostHardwareSerial: fmt.Sprintf("serial%d", i),
		})
	}

	ds.GetHostCalendarEventByEmailFunc = func(ctx context.Context, email string) (*mobius.HostCalendarEvent, *mobius.CalendarEvent, error) {
		hostID, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(email, "user"), "@example.com"))
		require.NoError(t, err)
		if hostID%2 == 0 {
			return nil, nil, notFoundErr{}
		}
		require.Contains(t, eventPerHost, uint(hostID)) //nolint:gosec // dismiss G115
		return &mobius.HostCalendarEvent{
			ID:              uint(hostID), //nolint:gosec // dismiss G115
			HostID:          uint(hostID), //nolint:gosec // dismiss G115
			CalendarEventID: uint(hostID), //nolint:gosec // dismiss G115
			WebhookStatus:   mobius.CalendarWebhookStatusNone,
		}, eventPerHost[uint(hostID)], nil //nolint:gosec // dismiss G115
	}

	ds.DeleteCalendarEventFunc = func(ctx context.Context, calendarEventID uint) error {
		return nil
	}

	err = cronCalendarEvents(ctx, ds, distributedLock, defaultCalendarConfig, logger)
	require.NoError(t, err)

	createdCalendarEvents = calendar.ListGoogleMockEvents()
	require.Len(t, createdCalendarEvents, 0)
}

// TestEventBody tests generation of the event body.
func TestEventBody(t *testing.T) {
	ds := new(mock.Store)
	ctx := context.Background()
	logger := kitlog.With(kitlog.NewLogfmtLogger(os.Stdout))
	t.Cleanup(
		func() {
			calendar.ClearMockEvents()
			calendar.ClearMockChannels()
		},
	)

	//
	// Test setup
	//
	//	team1:
	//
	//	policyID1 (calendar) -- has description and resolution
	//	policyID2 (calendar) -- has description, but blank resolution
	//	policyID3 (calendar) -- has description, but nil resolution
	//	policyID4 (calendar) -- has no description, but has resolution
	//  policyID5 (calendar) -- returns error on lookup
	//
	// 	hostID1 not passing policyID1
	// 	hostID2 not passing policyID2
	// 	hostID3 not passing policyID3
	// 	hostID4 not passing policyID4
	//  hostID5 not passing policies 1,2,3,4
	//  hostID6 also not passing policyID1
	//  hostID7 not passing policyID5
	//

	const orgName = "Test Organization"
	ds.AppConfigFunc = func(ctx context.Context) (*mobius.AppConfig, error) {
		return &mobius.AppConfig{
			OrgInfo: mobius.OrgInfo{
				OrgName: orgName,
			},
			Integrations: mobius.Integrations{
				GoogleCalendar: []*mobius.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: map[string]string{
							mobius.GoogleCalendarEmail: "calendar-mock@example.com",
						},
					},
				},
			},
		}, nil
	}

	teamID1 := uint(1)
	ds.ListTeamsFunc = func(ctx context.Context, filter mobius.TeamFilter, opt mobius.ListOptions) ([]*mobius.Team, error) {
		return []*mobius.Team{
			{
				ID: teamID1,
				Config: mobius.TeamConfig{
					Integrations: mobius.TeamIntegrations{
						GoogleCalendar: &mobius.TeamGoogleCalendarIntegration{
							Enable:     true,
							WebhookURL: "https://foo.example.com",
						},
					},
				},
			},
		}, nil
	}

	policyID1 := uint(10)
	policyID2 := uint(11)
	policyID3 := uint(12)
	policyID4 := uint(13)
	policyID5 := uint(14)
	ds.GetCalendarPoliciesFunc = func(ctx context.Context, teamID uint) ([]mobius.PolicyCalendarData, error) {
		require.Equal(t, teamID1, teamID)
		return []mobius.PolicyCalendarData{
			{
				ID:   policyID1,
				Name: "Policy 1",
			},
			{
				ID:   policyID2,
				Name: "Policy 2",
			},
			{
				ID:   policyID3,
				Name: "Policy 3",
			},
			{
				ID:   policyID4,
				Name: "Policy 4",
			},
			{
				ID:   policyID5,
				Name: "Policy 5",
			},
		}, nil
	}

	hostID1, userEmail1, hostDisplayName1 := uint(100), "user1@example.com", "Host 1"
	hostID2, userEmail2, hostDisplayName2 := uint(101), "user2@example.com", "Host 2"
	hostID3, userEmail3, hostDisplayName3 := uint(102), "user3@example.com", "Host 3"
	hostID4, userEmail4, hostDisplayName4 := uint(103), "user4@example.com", "Host 4"
	hostID5, userEmail5, hostDisplayName5 := uint(104), "user5@example.com", "Host 5"
	hostID6, userEmail6, hostDisplayName6 := uint(105), "user6@example.com", "Host 6"
	hostID7, userEmail7, hostDisplayName7 := uint(106), "user7@example.com", "Host 7"

	ds.GetTeamHostsPolicyMembershipsFunc = func(
		ctx context.Context, domain string, teamID uint, policyIDs []uint, _ *uint,
	) ([]mobius.HostPolicyMembershipData, error) {
		require.Equal(t, "example.com", domain)
		require.Equal(t, teamID1, teamID)
		require.Equal(t, []uint{policyID1, policyID2, policyID3, policyID4, policyID5}, policyIDs)
		return []mobius.HostPolicyMembershipData{
			{
				HostID:           hostID1,
				Email:            userEmail1,
				HostDisplayName:  hostDisplayName1,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d", policyID1),
			},
			{
				HostID:           hostID2,
				Email:            userEmail2,
				HostDisplayName:  hostDisplayName2,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d", policyID2),
			},
			{
				HostID:           hostID3,
				Email:            userEmail3,
				HostDisplayName:  hostDisplayName3,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d", policyID3),
			},
			{
				HostID:           hostID4,
				Email:            userEmail4,
				HostDisplayName:  hostDisplayName4,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d", policyID4),
			},
			{
				HostID:           hostID5,
				Email:            userEmail5,
				HostDisplayName:  hostDisplayName5,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d,%d,%d,%d", policyID1, policyID2, policyID3, policyID4),
			},
			{
				HostID:           hostID6,
				Email:            userEmail6,
				HostDisplayName:  hostDisplayName6,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d", policyID1),
			},
			{
				HostID:           hostID7,
				Email:            userEmail7,
				HostDisplayName:  hostDisplayName7,
				Passing:          false,
				FailingPolicyIDs: fmt.Sprintf("%d", policyID5),
			},
		}, nil
	}
	ds.PolicyLiteFunc = func(ctx context.Context, policyID uint) (*mobius.PolicyLite, error) {
		switch policyID {
		case policyID1:
			return &mobius.PolicyLite{
				ID:          policyID1,
				Description: "Description for policy 1",
				Resolution:  ptr.String("Resolution for policy 1"),
			}, nil
		case policyID2:
			return &mobius.PolicyLite{
				ID:          policyID2,
				Description: "Description for policy 2",
				Resolution:  ptr.String(""),
			}, nil
		case policyID3:
			return &mobius.PolicyLite{
				ID:          policyID2,
				Description: "Description for policy 3",
				Resolution:  nil,
			}, nil
		case policyID4:
			return &mobius.PolicyLite{
				ID:         policyID4,
				Resolution: ptr.String("Resolution for policy 4"),
			}, nil
		case policyID5:
			return nil, notFoundErr{}
		default:
			t.Errorf("unexpected policy ID: %d", policyID)
			return nil, nil
		}
	}

	ds.GetHostCalendarEventByEmailFunc = func(ctx context.Context, email string) (*mobius.HostCalendarEvent, *mobius.CalendarEvent, error) {
		return nil, nil, notFoundErr{}
	}

	var eventsMu sync.Mutex
	calendarEvents := make(map[uint]*mobius.CalendarEvent)
	hostCalendarEvents := make(map[uint]*mobius.HostCalendarEvent)

	ds.CreateOrUpdateCalendarEventFunc = func(
		ctx context.Context,
		uuid string,
		email string,
		startTime, endTime time.Time,
		data []byte,
		timeZone *string,
		hostID uint,
		webhookStatus mobius.CalendarWebhookStatus,
	) (*mobius.CalendarEvent, error) {
		assert.NotEmpty(t, uuid)
		require.Equal(t, mobius.CalendarWebhookStatusNone, webhookStatus)
		require.NotEmpty(t, data)
		require.NotZero(t, startTime)
		require.NotZero(t, endTime)

		eventsMu.Lock()
		calendarEventID := uint(len(calendarEvents) + 1) //nolint:gosec // dismiss G115
		calendarEvents[hostID] = &mobius.CalendarEvent{
			ID:        calendarEventID,
			Email:     email,
			StartTime: startTime,
			EndTime:   endTime,
			Data:      data,
		}
		hostCalendarEventID := uint(len(hostCalendarEvents) + 1) //nolint:gosec // dismiss G115
		hostCalendarEvents[hostID] = &mobius.HostCalendarEvent{
			ID:              hostCalendarEventID,
			HostID:          hostID,
			CalendarEventID: calendarEventID,
			WebhookStatus:   webhookStatus,
		}
		eventsMu.Unlock()
		return nil, nil
	}

	pool := redistest.SetupRedis(t, t.Name(), false, false, false)
	err := cronCalendarEvents(ctx, ds, redis_lock.NewLock(pool), defaultCalendarConfig, logger)
	require.NoError(t, err)

	numberOfEvents := 7
	eventsMu.Lock()
	require.Len(t, calendarEvents, numberOfEvents)
	require.Len(t, hostCalendarEvents, numberOfEvents)
	eventsMu.Unlock()

	createdCalendarEvents := calendar.ListGoogleMockEvents()
	require.Len(t, createdCalendarEvents, numberOfEvents)
	for _, hostCalEvent := range hostCalendarEvents {
		var details map[string]string
		err = json.Unmarshal(calendarEvents[hostCalEvent.HostID].Data, &details)
		require.NoError(t, err)
		// What Google Calendar calls the "Description" is what Mobius calls the "Body," since the Body
		// contains a description and a resolution.
		eventBody := createdCalendarEvents[details["id"]].Description
		switch hostCalEvent.HostID {
		case hostID1:
			assert.Contains(t, eventBody, fmt.Sprintf(`%s %s (Host 1).`, orgName, mobius.CalendarBodyStaticHeader))
			assert.Contains(t, eventBody, "Description for policy 1")
			assert.Contains(t, eventBody, "Resolution for policy 1")
		case hostID6:
			assert.Contains(t, eventBody, fmt.Sprintf(`%s %s (Host 6).`, orgName, mobius.CalendarBodyStaticHeader))
			assert.Contains(t, eventBody, "Description for policy 1")
			assert.Contains(t, eventBody, "Resolution for policy 1")
		default:
			assert.Contains(t, eventBody, fmt.Sprintf(`%s %s (Host`, orgName, mobius.CalendarBodyStaticHeader))
			assert.Contains(t, eventBody, mobius.CalendarDefaultResolution)
		}
	}
}
