CREATE EXTENSION citext;
CREATE DOMAIN domain_email AS citext
	CHECK ( VALUE ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );

CREATE DOMAIN domain_phone AS VARCHAR(11)
	CHECK ( VALUE ~ '\+\d{10}' );

CREATE TABLE orders (
  	order_uid VARCHAR(255) UNIQUE PRIMARY KEY,
  	track_number VARCHAR(255),
  	entry VARCHAR(255),
  	locale VARCHAR(4),
  	internal_signature TEXT,
  	customer_id VARCHAR(255),
  	delivery_service TEXT,
  	shardkey VARCHAR(255),
  	sm_id INTEGER,
  	date_created TIMESTAMP,
  	oof_shard VARCHAR(255)
);

CREATE TABLE delivery (
  	order_uid VARCHAR(255) UNIQUE PRIMARY KEY,
  	recipient_name VARCHAR(255),
  	phone domain_phone,
  	zip VARCHAR(255),
  	city VARCHAR(255),
  	address VARCHAR(255),
  	region VARCHAR(255),
  	email domain_email,
  	CONSTRAINT fk_order
  		FOREIGN KEY(order_uid)
  			REFERENCES orders(order_uid)
);

CREATE TABLE payments (
  	order_uid VARCHAR(255) UNIQUE PRIMARY KEY,
  	transaction VARCHAR(255),
  	request_id VARCHAR(255),
  	currency VARCHAR(255),
  	provider VARCHAR(255),
  	amount NUMERIC,
  	payment_dt INTEGER,
  	bank VARCHAR(255),
  	delivery_cost NUMERIC,
  	goods_total INTEGER,
	custom_fee NUMERIC,
	CONSTRAINT fk_order
  		FOREIGN KEY(order_uid)
  			REFERENCES orders(order_uid)
);

CREATE TABLE items (
	chrt_id INTEGER UNIQUE PRIMARY KEY,
	order_uid VARCHAR(255),
  	track_number VARCHAR(255),
  	price INTEGER,
  	rid VARCHAR(255),
  	chrt_name VARCHAR(255),
  	sale INTEGER,
  	chrt_size VARCHAR(255),
  	total_price NUMERIC,
  	nm_id INTEGER,
  	brand VARCHAR(255),
  	status INTEGER,
  	CONSTRAINT fk_order
  		FOREIGN KEY(order_uid)
  			REFERENCES orders(order_uid)
);