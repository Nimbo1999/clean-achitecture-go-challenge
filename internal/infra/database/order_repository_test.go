package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/nimbo1999/20-CleanArch/internal/entity"
	"github.com/stretchr/testify/suite"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)
	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestGivenTwoOrderCreated_WhenListed_ThenShouldDisplayTwoOrders() {
	// Creating orders
	order, err := entity.NewOrder("first", 10.0, 0.5)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	secondOrder, err := entity.NewOrder("second", 100.0, 0.5)
	suite.NoError(err)
	suite.NoError(secondOrder.CalculateFinalPrice())
	// Saving orders
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)
	err = repo.Save(secondOrder)
	suite.NoError(err)
	// Retrieving orders
	ctx := context.Background()
	orders, err := repo.List(ctx)
	suite.NoError(err)
	suite.NotEmpty(orders)
	suite.Len(orders, 2)
	suite.Contains(orders, order)
	suite.Contains(orders, secondOrder)
}
