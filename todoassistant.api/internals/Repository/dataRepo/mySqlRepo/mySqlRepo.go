package mySqlRepo

import (
	"context"
	"database/sql"
	"log"

	"fmt"
	"test-va/internals/Repository/dataRepo"
	"test-va/internals/entity/dataEntity"
)

type sqlRepo struct {
	conn *sql.DB
}

func (s *sqlRepo) GetCountries(ctx context.Context) ([]*dataEntity.Country, error) {
	stmt := fmt.Sprintf(`
		SELECT id, name, tel_code FROM Countries
	`)

	rows, err := s.conn.QueryContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var countries []*dataEntity.Country

	for rows.Next() {
		var country dataEntity.Country
		err := rows.Scan(
			&country.CountryId,
			&country.CountryName,
			&country.TelCode,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		countries = append(countries, &country)
	}
	if rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

func NewDataSqlRepo(conn *sql.DB) dataRepo.DataRepository {
	return &sqlRepo{conn: conn}
}
