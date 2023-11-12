package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/juanjoss/off-users-service/ports"
	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository() *userRepository {
	source := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)

	db, err := sqlx.Connect(os.Getenv("DB_DRIVER"), source)
	if err != nil {
		log.Fatalf("unable to connect to DB: %v", err.Error())
	}

	repo := &userRepository{
		db: db,
	}

	return repo
}

func (ur *userRepository) Register(request ports.RegisterRequest) error {
	var id int

	row := ur.db.QueryRow(
		`INSERT INTO users (first_name, last_name, email, password) 
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		request.User.FirstName,
		request.User.LastName,
		request.User.Email,
		request.User.Password,
	)
	if err := row.Scan(&id); err != nil {
		return err
	}

	for _, ssd := range request.SSDs {
		ssd.UserId = id
		_, err := ur.db.NamedExec(
			`INSERT INTO ssds (user_id, mac_address)
			VALUES (:user_id, :mac_address)`,
			ssd,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ur *userRepository) AddProductToSSD(ssdId int, barcode string, quantity int) error {
	_, err := ur.db.Exec(
		`INSERT INTO product_ssds (ssd_id, barcode, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`,
		ssdId, barcode, quantity,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) RandomSSDId() (int, error) {
	var id int

	err := ur.db.Get(&id, `SELECT id FROM ssds ORDER BY RANDOM() LIMIT 1`)
	if err != nil {
		return id, err
	}

	return id, nil
}
