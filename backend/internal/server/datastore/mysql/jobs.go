package mysql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/notawar/mobius/internal/server/contexts/ctxerr"
	"github.com/notawar/mobius/internal/server/mobius"
)

func (ds *Datastore) NewJob(ctx context.Context, job *mobius.Job) (*mobius.Job, error) {
	query := `
INSERT INTO jobs (
    name,
    args,
    state,
    retries,
    error,
    not_before
)
VALUES (?, ?, ?, ?, ?, COALESCE(?, NOW()))
`
	var notBefore *time.Time
	if !job.NotBefore.IsZero() {
		notBefore = &job.NotBefore
	}
	result, err := ds.writer(ctx).ExecContext(ctx, query, job.Name, job.Args, job.State, job.Retries, job.Error, notBefore)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	job.ID = uint(id) //nolint:gosec // dismiss G115

	return job, nil
}

func (ds *Datastore) GetQueuedJobs(ctx context.Context, maxNumJobs int, now time.Time) ([]*mobius.Job, error) {
	query := `
SELECT
    id, created_at, updated_at, name, args, state, retries, error, not_before
FROM
    jobs
WHERE
    state = ? AND
    not_before <= ?
ORDER BY
    updated_at ASC
LIMIT ?
`

	if now.IsZero() {
		now = time.Now().UTC()
	}

	var jobs []*mobius.Job
	err := sqlx.SelectContext(ctx, ds.reader(ctx), &jobs, query, mobius.JobStateQueued, now, maxNumJobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (ds *Datastore) UpdateJob(ctx context.Context, id uint, job *mobius.Job) (*mobius.Job, error) {
	query := `
UPDATE jobs
SET
    state = ?,
    retries = ?,
    error = ?,
    not_before = COALESCE(?, NOW())
WHERE
    id = ?
`
	var notBefore *time.Time
	if !job.NotBefore.IsZero() {
		notBefore = &job.NotBefore
	}
	_, err := ds.writer(ctx).ExecContext(ctx, query, job.State, job.Retries, job.Error, notBefore, id)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (ds *Datastore) CleanupWorkerJobs(ctx context.Context, failedSince, completedSince time.Duration) (int64, error) {
	// using not_before instead of created_at/updated_at to be able to use the
	// existing index, and the difference between those timestamps will be
	// minimal (max 5 retries for failed jobs, with a few hours difference).
	const stmt = `
	DELETE FROM
		jobs
	WHERE
		(state = ? AND not_before < ?) OR
		(state = ? AND not_before < ?)
`

	now := time.Now().UTC()
	failedBefore := now.Add(-failedSince)
	completedBefore := now.Add(-completedSince)

	res, err := ds.writer(ctx).ExecContext(ctx, stmt,
		mobius.JobStateFailure, failedBefore,
		mobius.JobStateSuccess, completedBefore)
	if err != nil {
		return 0, ctxerr.Wrap(ctx, err, "cleanup worker jobs")
	}
	n, _ := res.RowsAffected()
	return n, nil
}
