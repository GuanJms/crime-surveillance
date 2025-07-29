package data

import (
	"context"
	"errors"
	patroldb "patrolServiceApp/data/gen"
	"patrolServiceApp/ptr"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbTimeout = time.Second * 3
)

var ErrNoContent error = errors.New("no content error")

var ErrNotFound error = errors.New("not found error")

var ErrAlreadyExists error = errors.New("already exists error")

type PostgresRepository struct {
	q *patroldb.Queries
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		q: patroldb.New(pool),
	}
}

func (repo *PostgresRepository) InsertPatrol(p *Patrol) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var street, city, state *string
	var latitude, longitude *float64

	if p.Location != nil {
		street = p.Location.Street
		city = p.Location.City
		state = p.Location.State
		latitude = p.Location.Latitude
		longitude = p.Location.Longitude
	}

	params := patroldb.InsertPatrolParams{
		UserID:      parseUUID(*p.UserID),
		OfficerID:   ptr.Deref(p.OfficerID),
		OfficerName: ptr.Deref(p.OfficerName),
		Status:      ptr.Deref((*patroldb.PatrolStatus)(p.Status)),
		Street:      wrapString(street),
		City:        wrapString(city),
		State:       wrapString(state),
		Latitude:    wrapFloat(latitude),
		Longitude:   wrapFloat(longitude),
		CreatedAt:   wrapTime(p.CreatedAt),
		UpdatedAt:   wrapTime(p.UpdatedAt),
	}

	id, err := repo.q.InsertPatrol(ctx, params)
	if err != nil {
		// if error is 23505, then the patrol already exists
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" {
				return "", ErrAlreadyExists
			}
		}
		return "", err
	}

	return id.String(), nil
}

func (repo *PostgresRepository) GetAllPatrols(*GetPatrolInfoRequest) ([]*Patrol, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// log.Println("GetAllPatrols - Receiving get all patrols request")

	patrolProfiles, err := repo.q.GetAllPatrols(ctx)
	if err != nil {
		return nil, err
	}

	// log.Printf("Receive patrol profiles: %v", patrolProfiles)

	patrols := make([]*Patrol, 0)

	for _, patrolProfile := range patrolProfiles {
		p := toPatrol(patrolProfile)
		// log.Printf("GetAllPatrols - Converted patrol: %v", p.String())
		patrols = append(patrols, p)
	}

	return patrols, nil
}

func (repo *PostgresRepository) PutPatrolInfo(req *UpdatePatrolInfoRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	params := patroldb.UpdateAllPatrolByUserIDParams{
		OfficerID:   ptr.Deref(req.OfficerId),
		OfficerName: ptr.Deref(req.OfficerName),
		Status:      ptr.Deref((*patroldb.PatrolStatus)(req.Status)),
		Street:      wrapString(req.Street),
		City:        wrapString(req.City),
		State:       wrapString(req.State),
		Latitude:    wrapFloat(req.Latitude),
		Longitude:   wrapFloat(req.Longitude),
		UserID:      parseUUID(req.UserId),
	}

	rows, err := repo.q.UpdateAllPatrolByUserID(ctx, params)
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (repo *PostgresRepository) PatchPatrolInfo(req *UpdatePatrolInfoRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var status patroldb.NullPatrolStatus
	if req.Status != nil {
		status = patroldb.NullPatrolStatus{
			PatrolStatus: patroldb.PatrolStatus(ptr.Deref(req.Status)),
			Valid:        true,
		}
	}

	params := patroldb.IgnoreNullUpdatePatrolByUserIDParams{
		OfficerID:   wrapString(req.OfficerId),
		OfficerName: wrapString(req.OfficerName),
		Status:      status,
		Street:      wrapString(req.Street),
		City:        wrapString(req.City),
		State:       wrapString(req.State),
		Latitude:    wrapFloat(req.Latitude),
		Longitude:   wrapFloat(req.Longitude),
		UserID:      parseUUID(req.UserId),
	}

	rows, err := repo.q.IgnoreNullUpdatePatrolByUserID(ctx, params)
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func toPatrol(src patroldb.PatrolProfile) *Patrol {
	p := &Patrol{
		UserID:      ptr.Of(src.UserID.String()),
		OfficerID:   ptr.Of(src.OfficerID),
		OfficerName: ptr.Of(src.OfficerName),
		Status:      ptr.Of(string(src.Status)),
		CreatedAt:   timePtr(&src.CreatedAt),
		UpdatedAt:   timePtr(&src.UpdatedAt),
		Location: &Location{
			Street:    textPtr(&src.Street),
			City:      textPtr(&src.City),
			State:     textPtr(&src.State),
			Latitude:  floatPtr(&src.Latitude),
			Longitude: floatPtr(&src.Longitude),
		},
	}
	return p
}

func floatPtr(f *pgtype.Float8) *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

func timePtr(t *pgtype.Timestamp) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func textPtr(t *pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func parseUUID(s string) pgtype.UUID {
	var id pgtype.UUID
	_ = id.Scan(s)
	return id
}

func wrapTime(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}
func wrapString(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}

func wrapFloat(f *float64) pgtype.Float8 {
	if f != nil {
		return pgtype.Float8{Float64: *f, Valid: true}
	}
	return pgtype.Float8{Valid: false}
}
