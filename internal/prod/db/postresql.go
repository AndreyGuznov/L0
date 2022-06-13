package db

import (
	"context"
	"fmt"
	"log"

	"example.com/m/internal/prod/models"
	"example.com/m/internal/prod/service"
	"example.com/m/pkg/client/postgre"
	"github.com/jackc/pgconn"
)

const (
	forInsertProductSQL  = "INSERT INTO product (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING product_id, track_number;"
	forInsertDeliverySQL = "INSERT INTO delivery (id, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING name;"
	forInsertPaymentSQL  = "INSERT INTO payment (id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;"
	forInsertItemSQL     = "INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING track_number;"
	findProductByIdSQL   = "SELECT product.order_uid, product.track_number, product.entry, product.locale, delivery.name, delivery.phone, delivery.zip, delivery.city, delivery.address, delivery.region, delivery.email, payment.transaction, payment.request_id, payment.currency, payment.provider, payment.amount, payment.payment_dt, payment.bank, payment.delivery_cost, payment.goods_total, payment.custom_fee, product.internal_signature, product.customer_id, product.delivery_service, product.shardkey, product.sm_id, product.date_created, product.oof_shard FROM product, delivery, payment WHERE product_id=$1 AND delivery.id=$1 AND payment.id=$1"
)

type repos struct {
	client postgre.Client
}

func (r *repos) InsertProd(ctx context.Context, product *models.Product) error {
	if err := r.client.QueryRow(ctx, forInsertProductSQL,
		product.Uid,
		product.TrackNum,
		product.Entry,
		product.Locale,
		product.Signature,
		product.Customer,
		product.DeliveryServ,
		product.Shardkey,
		product.SmId,
		product.DateOf,
		product.OofShard).Scan(&product.Product_id, &product.TrackNum); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			msgEr := fmt.Sprintf("Err off query1: %s", pgErr.Message)
			log.Println(msgEr)
			return err
		}
		return err
	}
	if err := r.client.QueryRow(ctx, forInsertDeliverySQL,
		product.Product_id,
		product.Delivery.Name,
		product.Delivery.Phone,
		product.Delivery.Zip,
		product.Delivery.City,
		product.Delivery.Address,
		product.Delivery.Region,
		product.Delivery.Email).Scan(&product.Delivery.Name); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			msgEr := fmt.Sprintf("Err off query2: %s", pgErr.Message)
			log.Println(msgEr)
			return err
		}
		return err
	}
	if err := r.client.QueryRow(ctx, forInsertPaymentSQL,
		product.Product_id,
		product.Payment.Transaction,
		product.Payment.Request,
		product.Payment.Currency,
		product.Payment.Provider,
		product.Payment.Amount,
		product.Payment.Paymen,
		product.Payment.Bank,
		product.Payment.Deliver,
		product.Payment.Goods,
		product.Payment.Custom).Scan(&product.Payment.Id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			msgEr := fmt.Sprintf("Err off query3: %s", pgErr.Message)
			log.Println(msgEr)
			return err
		}
		return err
	}

	for i := len(product.Items); i > 0; i-- {
		if err := r.client.QueryRow(ctx, forInsertItemSQL,
			product.Items[len(product.Items)-i].Chrt,
			product.Items[len(product.Items)-i].Number,
			product.Items[len(product.Items)-i].Price,
			product.Items[len(product.Items)-i].Rid,
			product.Items[len(product.Items)-i].NameOf,
			product.Items[len(product.Items)-i].Sale,
			product.Items[len(product.Items)-i].Size,
			product.Items[len(product.Items)-i].Total,
			product.Items[len(product.Items)-i].Nm,
			product.Items[len(product.Items)-i].Brand,
			product.Items[len(product.Items)-i].Status).Scan(&product.Items[len(product.Items)-i].Number); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				msgEr := fmt.Sprintf("Err off query4: %s", pgErr.Message)
				log.Println(msgEr)
				return err
			}
			return err
		}
	}

	return nil
}

func (r *repos) FindByIdProd(ctx context.Context, product_id int) (models.Product, error) {
	var product models.Product
	err := r.client.QueryRow(ctx, findProductByIdSQL, product_id).Scan(
		&product.Uid,
		&product.TrackNum,
		&product.Entry,
		&product.Locale,
		&product.Delivery.Name,
		&product.Delivery.Phone,
		&product.Delivery.Zip,
		&product.Delivery.City,
		&product.Delivery.Address,
		&product.Delivery.Region,
		&product.Delivery.Email,
		&product.Payment.Transaction,
		&product.Payment.Request,
		&product.Payment.Currency,
		&product.Payment.Provider,
		&product.Payment.Amount,
		&product.Payment.Paymen,
		&product.Payment.Bank,
		&product.Payment.Deliver,
		&product.Payment.Goods,
		&product.Payment.Custom,
		&product.Signature,
		&product.Customer,
		&product.DeliveryServ,
		&product.Shardkey,
		&product.SmId,
		&product.DateOf,
		&product.OofShard)
	if err != nil {
		log.Println(err)
	}

	queryWithTrack := fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE items.track_number in ('%s')", product.TrackNum)

	rows, err := r.client.Query(ctx, queryWithTrack)
	if err != nil {
		log.Println("Err of geting Items from DB")
	}
	items := make([]models.Items, 0)

	for rows.Next() {
		var itm models.Items
		// for i := len(product.Items); i > 0; i-- {

		err = rows.Scan(
			&itm.Chrt,
			&itm.Number,
			&itm.Price,
			&itm.NameOf,
			&itm.Rid,
			&itm.Sale,
			&itm.Size,
			&itm.Total,
			&itm.Nm,
			&itm.Brand,
			&itm.Status)

		if err != nil {
			return product, err
		}
		// }
		items = append(items, itm)
	}

	if err = rows.Err(); err != nil {
		return product, err
	}
	product.Items = append(product.Items, items...)
	return product, nil
}

func (r *repos) GetAllDataDB(ctx context.Context) ([]models.Product, error) {

	allData := make([]models.Product, 0)

	queryGetAllProducts := "SELECT product_id, order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM product"
	var product models.Product
	rows, err := r.client.Query(ctx, queryGetAllProducts)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {

		err := rows.Scan(
			&product.Product_id,
			&product.Uid,
			&product.TrackNum,
			&product.Entry,
			&product.Locale,
			&product.Signature,
			&product.Customer,
			&product.DeliveryServ,
			&product.Shardkey,
			&product.SmId,
			&product.DateOf,
			&product.OofShard)
		if err != nil {
			log.Println(err)
		}
		queryGetAllDelivery := "SELECT name, phone, zip, city, address, region, email FROM delivery, product WHERE product.product_id=delivery.id"
		var delivery models.Delivery
		err = r.client.QueryRow(ctx, queryGetAllDelivery).Scan(
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email)
		if err != nil {
			log.Println(err)
		}
		product.Delivery = delivery
		queryGetAllPayment := "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment, product WHERE product.product_id=payment.id"
		var payment models.Payment
		err = r.client.QueryRow(ctx, queryGetAllPayment).Scan(
			&payment.Transaction,
			&payment.Request,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.Paymen,
			&payment.Bank,
			&payment.Deliver,
			&payment.Goods,
			&payment.Custom)
		if err != nil {
			log.Println(err)
		}
		product.Payment = payment
		queryWithTrack := fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE items.track_number in ('%s')", product.TrackNum)

		rows, err := r.client.Query(ctx, queryWithTrack)
		if err != nil {
			log.Println("Err of geting Items from DB")
		}
		items := make([]models.Items, 0)

		for rows.Next() {
			var itm models.Items
			// for i := len(product.Items); i > 0; i-- {

			err = rows.Scan(
				&itm.Chrt,
				&itm.Number,
				&itm.Price,
				&itm.NameOf,
				&itm.Rid,
				&itm.Sale,
				&itm.Size,
				&itm.Total,
				&itm.Nm,
				&itm.Brand,
				&itm.Status)

			if err != nil {
				return nil, err
			}
			// }
			items = append(items, itm)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
		product.Items = append(product.Items, items...)

		// queryGetAll := "SELECT product.product_id, product.order_uid, product.track_number, product.entry, product.locale, delivery.name, delivery.phone, delivery.zip, delivery.city, delivery.address, delivery.region, delivery.email, payment.transaction, payment.request_id, payment.currency, payment.provider, payment.amount, payment.payment_dt, payment.bank, payment.delivery_cost, payment.goods_total, payment.custom_fee, product.internal_signature, product.customer_id, product.delivery_service, product.shardkey, product.sm_id, product.date_created, product.oof_shard FROM product, delivery, payment"
		// var product models.Product
		// rows, err := r.client.Query(ctx, queryGetAll)
		// if err != nil {
		// 	log.Println(err)
		// }
		// for rows.Next() {

		// 	err := rows.Scan(
		// 		&product.Product_id,
		// 		&product.Uid,
		// 		&product.TrackNum,
		// 		&product.Entry,
		// 		&product.Locale,
		// 		&product.Delivery.Name,
		// 		&product.Delivery.Phone,
		// 		&product.Delivery.Zip,
		// 		&product.Delivery.City,
		// 		&product.Delivery.Address,
		// 		&product.Delivery.Region,
		// 		&product.Delivery.Email,
		// 		&product.Payment.Transaction,
		// 		&product.Payment.Request,
		// 		&product.Payment.Currency,
		// 		&product.Payment.Provider,
		// 		&product.Payment.Amount,
		// 		&product.Payment.Paymen,
		// 		&product.Payment.Bank,
		// 		&product.Payment.Deliver,
		// 		&product.Payment.Goods,
		// 		&product.Payment.Custom,
		// 		&product.Signature,
		// 		&product.Customer,
		// 		&product.DeliveryServ,
		// 		&product.Shardkey,
		// 		&product.SmId,
		// 		&product.DateOf,
		// 		&product.OofShard)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// queryWithTrack := "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items;"

		// rows, err := r.client.Query(ctx, queryWithTrack)
		// if err != nil {
		// 	log.Println(err)
		// }
		// // items := make([]models.Items, 0)

		// for rows.Next() {
		// 	var itm models.Items
		// 	err := rows.Scan(
		// 		&itm.Chrt,
		// 		&itm.Number,
		// 		&itm.Price,
		// 		&itm.NameOf,
		// 		&itm.Rid,
		// 		&itm.Sale,
		// 		&itm.Size,
		// 		&itm.Total,
		// 		&itm.Nm,
		// 		&itm.Brand,
		// 		&itm.Status)
		// 	// for rows.Next() {

		// 	// 	err := rows.Scan(
		// 	// 		&itm.Chrt,
		// 	// 		&itm.Number,
		// 	// 		&itm.Price,
		// 	// 		&itm.NameOf,
		// 	// 		&itm.Rid,
		// 	// 		&itm.Sale,
		// 	// 		&itm.Size,
		// 	// 		&itm.Total,
		// 	// 		&itm.Nm,
		// 	// 		&itm.Brand,
		// 	// 		&itm.Status)

		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	for i := len(product.Items); i > 0; i-- {
		// 		if product.Items[i].Number == product.TrackNum {
		// 			product.Items = append(product.Items, itm)
		// 		}
		// 	}

		// }
		// // if err = rows.Err(); err != nil {
		// // 	log.Println(err)
		// // }

		// // if err != nil {
		// // 	log.Println(err)
		// // }
		allData = append(allData, product)
	}

	return allData, nil
}

func NewRepos(client postgre.Client) service.Storage {

	return &repos{
		client: client,
	}
}

// func NewRepos() service.Storage {

// 	PostgeSqlClient, err := postgre.NewClient(context.Background(), 5)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return &repos{
// 		client: PostgeSqlClient,
// 	}
// } 											//  При каждом вызове подключение к базе
