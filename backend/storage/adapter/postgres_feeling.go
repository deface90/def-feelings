package adapter

import (
	"context"
	"fmt"
	"github.com/deface90/def-feelings/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func (p *Postgres) CreateStatus(ctx context.Context, status *storage.Status) (int64, error) {
	var feelingsIDs []int64
	if len(status.Feelings) == 0 {
		return 0, errors.New("Failed to create status: feelings list is empty (probably, unvalidated)")
	}

	for _, f := range status.Feelings {
		feeling, err := p.getFeeling(ctx, f)
		if err == pgx.ErrNoRows {
			id, err := p.createFeeling(ctx, &storage.Feeling{Title: f})
			if err != nil {
				return 0, err
			}
			feelingsIDs = append(feelingsIDs, id)
			continue
		} else if err != nil {
			return 0, errors.Wrap(err, "Failed to get feeling")
		}
		feelingsIDs = append(feelingsIDs, feeling.ID)
	}

	query := fmt.Sprintf(`INSERT INTO %s (user_id, message, created) VALUES ($1, $2, $3) RETURNING id`, postgresStatus)
	row := p.db.QueryRow(ctx, query, status.UserID, status.Message, status.Created.Format(time.RFC3339))
	var statusID int64
	err := row.Scan(&statusID)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to scan status")
	}

	query = fmt.Sprintf("INSERT INTO %s (status_id, feeling_id) VALUES ", postgresStatusFeeling)
	var values []string
	for _, fID := range feelingsIDs {
		values = append(values, fmt.Sprintf("(%v, %v)", statusID, fID))
	}
	query += strings.Join(values, ",")
	_, err = p.db.Exec(ctx, query)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to create feeling-status link")
	}

	return statusID, nil
}

func (p *Postgres) GetStatus(ctx context.Context, statusID int64) (*storage.Status, error) {
	query := fmt.Sprintf("SELECT s.* FROM %v s WHERE s.id=$1", postgresStatus)
	row := p.db.QueryRow(ctx, query, statusID)

	var status storage.Status
	err := row.Scan(&status.ID, &status.UserID, &status.Message, &status.Created)
	if err != nil {
		return nil, err
	}
	status.Feelings, err = p.getStatusFeelings(ctx, status.ID)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (p *Postgres) ListStatuses(ctx context.Context, request storage.ListStatusesRequest) ([]*storage.Status, int64, error) {
	if request.UserID == 0 {
		return nil, 0, errors.New("User ID is required")
	}
	baseQuery := fmt.Sprintf("SELECT DISTINCT s.* FROM %v s, %v f, %v sf WHERE s.id=sf.status_id AND f.id=sf.feeling_id AND s.user_id=$1",
		postgresStatus, postgresFeeling, postgresStatusFeeling)

	if len(request.Feelings) != 0 {
		baseQuery += fmt.Sprintf(" AND LOWER(f.title) IN %v", p.prepareInCondition(request.Feelings, true))
	}
	if !request.DatetimeStart.IsZero() {
		baseQuery += fmt.Sprintf(" AND created>='%v'", request.DatetimeStart.Format(time.RFC3339))
	}
	if !request.DatetimeEnd.IsZero() {
		baseQuery += fmt.Sprintf(" AND created<='%v'", request.DatetimeEnd.Format(time.RFC3339))
	}

	query := baseQuery + " ORDER BY created DESC LIMIT $2 OFFSET $3"
	rows, err := p.db.Query(ctx, query, request.UserID, request.Limit, request.Limit*request.Page)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to get status list")
	}

	var statuses []*storage.Status
	for rows.Next() {
		status := new(storage.Status)
		err = rows.Scan(&status.ID, &status.UserID, &status.Message, &status.Created)
		if err != nil {
			return nil, 0, err
		}
		status.Feelings, err = p.getStatusFeelings(ctx, status.ID)
		if err != nil {
			return nil, 0, err
		}
		statuses = append(statuses, status)
	}

	count, err := p.customCount(ctx, baseQuery, request.UserID)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to count statuses list")
	}

	return statuses, count, nil
}

func (p *Postgres) ListFeelings(ctx context.Context, request storage.ListFeelingsRequest) ([]*storage.Feeling, error) {
	query := fmt.Sprintf("SELECT * FROM %s", postgresFeeling)
	if request.Title != "" {
		query += fmt.Sprintf(" WHERE title ILIKE '%%%s%%'", p.escape(request.Title))
	}
	query += " LIMIT $1 OFFSET $2"

	rows, err := p.db.Query(ctx, query, request.Limit, request.Limit*request.Page)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get feelings list")
	}

	var feelings []*storage.Feeling
	for rows.Next() {
		f := new(storage.Feeling)
		err := rows.Scan(&f.ID, &f.Title)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to scan feeling")
		}

		feelings = append(feelings, f)
	}

	return feelings, nil
}

func (p *Postgres) GetFeelingsFrequency(ctx context.Context, request storage.FeelingsFrequencyRequest) ([]storage.FeelingFrequencyItem, error) {
	query := fmt.Sprintf(`
SELECT COUNT(sf.*) as feeling_count, sf.feeling_id, f.title FROM %s sf
INNER JOIN %s s ON sf.status_id=s.id`, postgresStatusFeeling, postgresStatus)
	if !request.DatetimeStart.IsZero() {
		query += fmt.Sprintf(" AND s.created>='%v'", request.DatetimeStart.Format(time.RFC3339))
	}
	if !request.DatetimeEnd.IsZero() {
		query += fmt.Sprintf(" AND s.created<='%v'", request.DatetimeEnd.Format(time.RFC3339))
	}

	query += fmt.Sprintf(` LEFT JOIN %s f ON sf.feeling_id=f.id 
WHERE sf.status_id=s.id AND sf.feeling_id=f.id AND s.user_id=$1
GROUP BY sf.feeling_id, f.title
ORDER BY feeling_count DESC
LIMIT $2 OFFSET $3`, postgresFeeling)

	rows, err := p.db.Query(ctx, query, request.UserID, request.Limit, request.Limit*request.Page)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get frequent feelings list")
	}

	var response []storage.FeelingFrequencyItem
	for rows.Next() {
		var item storage.FeelingFrequencyItem
		err := rows.Scan(&item.Frequency, &item.Feeling.ID, &item.Feeling.Title)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to scan frequent feeling")
		}

		response = append(response, item)
	}

	return response, nil
}

func (p *Postgres) createFeeling(ctx context.Context, feeling *storage.Feeling) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (title) VALUES ($1) RETURNING id`, postgresFeeling)
	row := p.db.QueryRow(ctx, query, feeling.Title)
	var id int64
	err := row.Scan(&id)
	return id, err
}

func (p *Postgres) getFeeling(ctx context.Context, title string) (*storage.Feeling, error) {
	query := fmt.Sprintf("SELECT * FROM %s u WHERE title=$1", postgresFeeling)
	row := p.db.QueryRow(ctx, query, title)

	var feeling storage.Feeling
	err := row.Scan(&feeling.ID, &feeling.Title)
	if err != nil {
		return nil, err
	}

	return &feeling, nil
}

func (p *Postgres) getStatusFeelings(ctx context.Context, statusID int64) ([]string, error) {
	query := fmt.Sprintf(`SELECT f.title FROM %v f, %v sf WHERE f.id=sf.feeling_id AND sf.status_id=$1`,
		postgresFeeling, postgresStatusFeeling)
	rows, err := p.db.Query(ctx, query, statusID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get feelings list by status")
	}

	var feelings []string
	for rows.Next() {
		var f string
		err := rows.Scan(&f)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to scan feeling title")
		}

		feelings = append(feelings, f)
	}

	return feelings, nil
}
