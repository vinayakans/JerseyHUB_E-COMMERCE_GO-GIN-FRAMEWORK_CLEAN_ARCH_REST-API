package repository

import (
	"fmt"
	"jerseyhub/pkg/domain"
	"jerseyhub/pkg/utils/models"

	"gorm.io/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (or *orderRepository) GetOrders(id int) ([]domain.Order, error) {

	var orders []domain.Order

	if err := or.DB.Raw("select * from orders where user_id=?", id).Scan(&orders).Error; err != nil {
		return []domain.Order{}, err
	}

	return orders, nil

}

func (ad *orderRepository) GetCart(id int) ([]models.GetCart, error) {

	var cart []models.GetCart

	if err := ad.DB.Raw("SELECT inventories.product_name,cart_products.quantity,cart_products.total_price AS Total FROM cart_products JOIN inventories ON cart_products.inventory_id=inventories.id WHERE user_id=$1", id).Scan(&cart).Error; err != nil {
		return []models.GetCart{}, err
	}
	fmt.Println(cart)
	return cart, nil

}

func (i *orderRepository) OrderItems(userid, addressid, paymentid int, total float64) (int, error) {

	var id int
	query := `
    INSERT INTO orders (user_id,address_id, payment_method_id, final_price)
    VALUES (?, ?, ?, ?)
    RETURNING id
    `
	i.DB.Raw(query, userid, addressid, paymentid, total).Scan(&id)

	return id, nil

}

func (i *orderRepository) AddOrderProducts(order_id int, cart []models.GetCart) error {

	query := `
    INSERT INTO order_items (order_id,inventory_id,quantity,total_price)
    VALUES (?, ?, ?, ?)
    `

	for _, v := range cart {
		var inv int
		if err := i.DB.Raw("select id from inventories where product_name=$1", v.ProductName).Scan(&inv).Error; err != nil {
			return err
		}

		if err := i.DB.Exec(query, order_id, inv, v.Quantity, v.Total).Error; err != nil {
			return err
		}
	}

	return nil

}

func (i *orderRepository) CancelOrder(status string, id int) error {

	if err := i.DB.Exec("update orders set order_status=$1 where id=$2", status, id).Error; err != nil {
		return err
	}

	return nil

}

func (i *orderRepository) EditOrderStatus(status string, id int) error {

	if err := i.DB.Exec("update orders set order_status=$1 where id=$2", status, id).Error; err != nil {
		return err
	}

	return nil

}

func (or *orderRepository) AdminOrders(status string) ([]domain.Order, error) {

	var orders []domain.Order
	if err := or.DB.Raw("select * from orders where order_status=$1", status).Scan(&orders).Error; err != nil {
		return []domain.Order{}, err
	}

	return orders, nil

}
