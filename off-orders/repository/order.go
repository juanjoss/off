package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/juanjoss/off-orders-service/model"
	_ "github.com/lib/pq"
)

type ordersRepository struct {
	db *sqlx.DB
}

func NewProductRepository() *ordersRepository {
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

	repo := &ordersRepository{
		db: db,
	}

	return repo
}

/*
Returns a random product.
*/
func (pr *ordersRepository) GetRandomProductBarcode() (string, error) {
	var response string

	err := pr.db.Get(&response, `SELECT barcode FROM products ORDER BY RANDOM() LIMIT 1`)
	if err != nil {
		return response, err
	}

	return response, nil
}

/*
Creates a product order.
*/
func (pr *ordersRepository) CreateProductOrder(order model.Order) (int, error) {
	var id int
	log.Println(order)
	rows, err := pr.db.Queryx(
		`INSERT INTO product_orders (ssd_id, product_id, timestamp, quantity, status) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING
		RETURNING id`,
		order.SsdId, order.ProductId, order.Timestamp, order.Quantity, order.Status,
	)
	if err != nil {
		return id, err
	}

	if rows.Next() {
		rows.Scan(&id)
	}

	return id, nil
}

/*
Updates a product order.
*/
func (pr *ordersRepository) UpdateProductOrder(order model.Order) error {
	_, err := pr.db.NamedExec(
		`UPDATE product_orders 
		SET status = :status
		WHERE id = :id`,
		order,
	)
	if err != nil {
		return err
	}

	return nil
}
