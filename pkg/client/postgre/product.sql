CREATE TABLE public.product (
    product_id serial primary KEY,
	order_uid VARCHAR(100) NULL,
	track_number VARCHAR(100) UNIQUE,
	entry VARCHAR(100) NULL,
	locale VARCHAR(100) NULL,
	internal_signature VARCHAR(100) NULL,
	customer_id VARCHAR(100) NULL,
	delivery_service VARCHAR(100) NULL,
	shardkey VARCHAR(100) NULL,
	sm_id INT NULL,
	date_created VARCHAR(100) NULL,
	oof_shard VARCHAR(100) NULL
);

CREATE TABLE public.delivery (
	id int not NULL,
	"name" VARCHAR(100) NULL,
	phone VARCHAR(100) NULL,
	zip VARCHAR(100) NULL,
	city VARCHAR(100) NULL,
	address VARCHAR(100) NULL,
	region VARCHAR(100) NULL,
	email VARCHAR(100) NULL,
    FOREIGN KEY (id) REFERENCES public.product (product_id) ON DELETE CASCADE
);

CREATE TABLE public.payment (
	id int not NULL,
	"transaction" VARCHAR(100) NULL,
	request_id VARCHAR(100) NULL,
	currency VARCHAR(100) NULL,
	provider VARCHAR(100) NULL,
	amount INT NULL,
	payment_dt BIGINT NULL,
	bank VARCHAR(100) NULL,
	delivery_cost INT NULL,
	goods_total INT NULL,
	custom_fee INT NULL,
    FOREIGN KEY (id) REFERENCES public.product (product_id) ON DELETE CASCADE
);

CREATE TABLE public.items (
	id serial not NULL,
	chrt_id BIGINT NULL,
	track_number VARCHAR(100) not NULL,
	price BIGINT NULL,
	rid VARCHAR(100) NULL,
	"name" VARCHAR(100) NULL,
	sale INT NULL,
	"size" VARCHAR(100) NULL,
	total_price BIGINT NULL,
	nm_id BIGINT NULL,
	brand VARCHAR(100) NULL,
	status INT NULL,
    FOREIGN KEY (track_number) REFERENCES public.product (track_number) ON DELETE CASCADE
);