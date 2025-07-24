package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(conn *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		Conn: conn,
	}
}

// TODO: Adding filtering conditions
func (repo *PostgresRepository) GetAllCrimes() ([]*Crime, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, reporter_id, description, status, street, city, state, 
	latitude, longitude, reported_at, created_at, updated_at 
	from crime order by reported_at
	`

	rows, err := repo.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crimes []*Crime

	for rows.Next() {
		var crime Crime
		err := rows.Scan(
			&crime.ID,
			&crime.ReporterID,
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
		crimes = append(crimes, &crime)
	}

	return crimes, nil
}

func (repo *PostgresRepository) InsertNewCrime(c *Crime) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO crime (
			id, reporter_id, description, status, 
			street, city, state, latitude, longitude
		)
		VALUES ($1, $2, $3, $4, 
				$5, $6, $7, $8, $9)
	`
	_, err := repo.Conn.ExecContext(ctx, query,
		c.ID,
		c.ReporterID,
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
