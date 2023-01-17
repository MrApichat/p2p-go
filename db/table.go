package db

var createTb = `
-- public.currencies definition

-- Drop table

-- DROP TABLE public.currencies;

CREATE TABLE public.currencies (
	id bigserial NOT NULL,
	"type" varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT currencies_pkey PRIMARY KEY (id)
);


-- public.failed_jobs definition

-- Drop table

-- DROP TABLE public.failed_jobs;

CREATE TABLE public.failed_jobs (
	id bigserial NOT NULL,
	uuid varchar(255) NOT NULL,
	"connection" text NOT NULL,
	queue text NOT NULL,
	payload text NOT NULL,
	"exception" text NOT NULL,
	failed_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT failed_jobs_pkey PRIMARY KEY (id),
	CONSTRAINT failed_jobs_uuid_unique UNIQUE (uuid)
);


-- public.migrations definition

-- Drop table

-- DROP TABLE public.migrations;

CREATE TABLE public.migrations (
	id serial4 NOT NULL,
	migration varchar(255) NOT NULL,
	batch int4 NOT NULL,
	CONSTRAINT migrations_pkey PRIMARY KEY (id)
);


-- public.password_resets definition

-- Drop table

-- DROP TABLE public.password_resets;

CREATE TABLE public.password_resets (
	email varchar(255) NOT NULL,
	"token" varchar(255) NOT NULL,
	created_at timestamp(0) NULL
);
CREATE INDEX password_resets_email_index ON public.password_resets USING btree (email);


-- public.payment_methods definition

-- Drop table

-- DROP TABLE public.payment_methods;

CREATE TABLE public.payment_methods (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT payment_methods_pkey PRIMARY KEY (id)
);


-- public.personal_access_tokens definition

-- Drop table

-- DROP TABLE public.personal_access_tokens;

CREATE TABLE public.personal_access_tokens (
	id bigserial NOT NULL,
	tokenable_type varchar(255) NOT NULL,
	tokenable_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	"token" varchar(64) NOT NULL,
	abilities text NULL,
	last_used_at timestamp(0) NULL,
	expires_at timestamp(0) NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT personal_access_tokens_pkey PRIMARY KEY (id),
	CONSTRAINT personal_access_tokens_token_unique UNIQUE (token)
);
CREATE INDEX personal_access_tokens_tokenable_type_tokenable_id_index ON public.personal_access_tokens USING btree (tokenable_type, tokenable_id);


-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	email_verified_at timestamp(0) NULL,
	"password" varchar(255) NOT NULL,
	remember_token varchar(100) NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT users_email_unique UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);


-- public.merchant_orders definition

-- Drop table

-- DROP TABLE public.merchant_orders;

CREATE TABLE public.merchant_orders (
	id bigserial NOT NULL,
	"type" varchar(255) NOT NULL,
	fiat_id int4 NOT NULL,
	coin_id int4 NOT NULL,
	merchant_id int4 NOT NULL,
	price float8 NOT NULL,
	available_coin float8 NOT NULL,
	lower_limit float8 NOT NULL,
	status varchar(255) NOT NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT merchant_orders_pkey PRIMARY KEY (id),
	CONSTRAINT merchant_orders_coin_id_foreign FOREIGN KEY (coin_id) REFERENCES public.currencies(id),
	CONSTRAINT merchant_orders_fiat_id_foreign FOREIGN KEY (fiat_id) REFERENCES public.currencies(id),
	CONSTRAINT merchant_orders_merchant_id_foreign FOREIGN KEY (merchant_id) REFERENCES public.users(id)
);


-- public.merchant_orders_payment_methods definition

-- Drop table

-- DROP TABLE public.merchant_orders_payment_methods;

CREATE TABLE public.merchant_orders_payment_methods (
	id bigserial NOT NULL,
	payment_method_id int4 NOT NULL,
	merchant_order_id int4 NOT NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT merchant_orders_payment_methods_pkey PRIMARY KEY (id),
	CONSTRAINT merchant_orders_payment_methods_merchant_order_id_foreign FOREIGN KEY (merchant_order_id) REFERENCES public.merchant_orders(id),
	CONSTRAINT merchant_orders_payment_methods_payment_method_id_foreign FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id)
);


-- public.trade_orders definition

-- Drop table

-- DROP TABLE public.trade_orders;

CREATE TABLE public.trade_orders (
	id bigserial NOT NULL,
	user_id int4 NOT NULL,
	merchant_order_id int4 NOT NULL,
	amount float8 NOT NULL,
	payment_method_id int4 NOT NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	total_price float8 NOT NULL,
	status varchar(255) NOT NULL,
	CONSTRAINT trade_orders_pkey PRIMARY KEY (id),
	CONSTRAINT trade_orders_merchant_order_id_foreign FOREIGN KEY (merchant_order_id) REFERENCES public.merchant_orders(id),
	CONSTRAINT trade_orders_payment_method_id_foreign FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id),
	CONSTRAINT trade_orders_user_id_foreign FOREIGN KEY (user_id) REFERENCES public.currencies(id)
);


-- public.transfer_orders definition

-- Drop table

-- DROP TABLE public.transfer_orders;

CREATE TABLE public.transfer_orders (
	id bigserial NOT NULL,
	coin_id int4 NOT NULL,
	sender_id int4 NOT NULL,
	receiver_id int4 NOT NULL,
	amount float8 NOT NULL,
	status varchar(255) NOT NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT transfer_orders_pkey PRIMARY KEY (id),
	CONSTRAINT transfer_orders_coin_id_foreign FOREIGN KEY (coin_id) REFERENCES public.currencies(id),
	CONSTRAINT transfer_orders_receiver_id_foreign FOREIGN KEY (receiver_id) REFERENCES public.users(id),
	CONSTRAINT transfer_orders_sender_id_foreign FOREIGN KEY (sender_id) REFERENCES public.users(id)
);


-- public.wallets definition

-- Drop table

-- DROP TABLE public.wallets;

CREATE TABLE public.wallets (
	id bigserial NOT NULL,
	user_id int4 NOT NULL,
	coin_id int4 NOT NULL,
	total float8 NOT NULL,
	in_order float8 NOT NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT wallets_pkey PRIMARY KEY (id),
	CONSTRAINT wallets_coin_id_foreign FOREIGN KEY (coin_id) REFERENCES public.currencies(id),
	CONSTRAINT wallets_user_id_foreign FOREIGN KEY (user_id) REFERENCES public.users(id)
);`
