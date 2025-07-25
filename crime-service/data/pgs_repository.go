package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const dbTimeout = time.Second * 3

var ErrNoConent error

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(conn *sql.DB) *PostgresRepository {
	ErrNoConent = errors.New("no content error")
	return &PostgresRepository{
		Conn: conn,
	}
}

// TODO: Adding filtering conditions
func (repo *PostgresRepository) GetAllCrimes() ([]*Crime, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, reporter_id, patrol_id, description, status, street, city, state, 
	latitude, longitude, reported_at, created_at, updated_at 
	from crime order by reported_at
	`

	rows, err := repo.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crimes []*Crime

	var patrolID sql.NullString

	for rows.Next() {
		var crime Crime
		err := rows.Scan(
			&crime.ID,
			&crime.ReporterID,
			&patrolID,
			&crime.Description,
			&crime.Status,
			&crime.Location.Street,
			&crime.Location.City,
			&crime.Location.State,
			&crime.Location.Latitude,
			&crime.Location.Longitude,
			&crime.ReportedAt,
			&crime.CreatedAt,
			&crime.UpdatedAt,
		)

		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		if patrolID.Valid {
			crime.PatrolID = patrolID.String
		}
		crimes = append(crimes, &crime)
	}

	return crimes, nil
}

func (repo *PostgresRepository) InsertNewCrime(c *Crime) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO crime (
			id, reporter_id, patrol_id, description, status, 
			street, city, state, latitude, longitude
		)
		VALUES ($1, $2, $3, $4, 
				$5, $6, $7, $8, $9, $10)
	`
	_, err := repo.Conn.ExecContext(ctx, query,
		c.ID,
		c.ReporterID,
		nilIfEmpty(c.PatrolID),
		c.Description,
		c.Status,
		c.Location.Street,
		c.Location.City,
		c.Location.State,
		c.Location.Latitude,
		c.Location.Longitude,
	)

	if err != nil {
		log.Printf("InsertNewCrime error: %v", err)
		return err
	}

	return nil
}

// Update the whole exisitng crime
func (repo *PostgresRepository) PutCrime(c *Crime) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE crime
		SET
			reporter_id = $1,
			patrol_id = $2,
			description = $3,
			status = $4,
			street = $5,
			city = $6,
			state = $7,
			latitude = $8,
			longitude = $9,
			updated_at = now()
		WHERE id = $10
		`

	resp, err := repo.Conn.ExecContext(ctx, query,
		c.ReporterID,
		nilIfEmpty(c.PatrolID),
		c.Description,
		c.Status,
		c.Location.Street,
		c.Location.City,
		c.Location.State,
		c.Location.Latitude,
		c.Location.Longitude,
		c.ID,
	)

	if row, err := resp.RowsAffected(); row == 0 || err != nil {
		if err != nil {
			return err
		} else {
			return ErrNoConent
		}
	}

	if err != nil {
		log.Printf("Updateing existing crime %s error: %v", c.ID, err)
		return err
	}

	return nil
}

// Partial udpate
func (repo *PostgresRepository) PatchCrime(update *CrimeUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query, args := buildingCrimeUpdateQuery(update)
	log.Println("Running ", query)

	resp, err := repo.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if row, err := resp.RowsAffected(); row == 0 || err != nil {
		if err != nil {
			log.Printf("Updateing existing crime %s error: %v", update.ID, err)
			return err
		} else {
			return ErrNoConent
		}
	}

	return nil
}

func buildingCrimeUpdateQuery(update *CrimeUpdate) (string, []any) {
	setParts := []string{}
	args := []any{}
	argPos := 1

	log.Println("Start building query")

	if update.ReporterID != nil {
		setParts = append(setParts, fmt.Sprintf("reporter_id = $%d", argPos))
		args = append(args, *update.ReporterID)
		argPos++
	}

	if update.PatrolID != nil {
		var newPatrolIDArg any
		setParts = append(setParts, fmt.Sprintf("patrol_id = $%d", argPos))
		if *update.PatrolID == "" {
			newPatrolIDArg = nil
		} else {
			newPatrolIDArg = *update.PatrolID
		}
		args = append(args, newPatrolIDArg)
		argPos++
	}

	if update.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *update.Description)
		argPos++
	}

	if update.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *update.Status)
		argPos++
	}

	if update.LocationUpdate != nil {

		if update.LocationUpdate.Street != nil {
			setParts = append(setParts, fmt.Sprintf("street = $%d", argPos))
			args = append(args, *update.LocationUpdate.Street)
			argPos++
		}

		if update.LocationUpdate.City != nil {
			setParts = append(setParts, fmt.Sprintf("city = $%d", argPos))
			args = append(args, *update.LocationUpdate.City)
			argPos++
		}

		if update.LocationUpdate.State != nil {
			setParts = append(setParts, fmt.Sprintf("state = $%d", argPos))
			args = append(args, *update.LocationUpdate.State)
			argPos++
		}

		if update.LocationUpdate.Latitude != nil {
			setParts = append(setParts, fmt.Sprintf("latitude = $%d", argPos))
			args = append(args, *update.LocationUpdate.Latitude)
			argPos++
		}

		if update.LocationUpdate.Longitude != nil {
			setParts = append(setParts, fmt.Sprintf("longitude = $%d", argPos))
			args = append(args, *update.LocationUpdate.Longitude)
			argPos++
		}
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = now()"))
	query := fmt.Sprintf("UPDATE crime SET %s WHERE id = $%d", strings.Join(setParts, ", "), argPos)
	args = append(args, update.ID)

	return query, args
}
