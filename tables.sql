CREATE TABLE public.bank_account (
    id character(36) NOT NULL,
    account_name character varying(256) NOT NULL,
    account_number character varying(256) NOT NULL,
    bank character varying(256) NOT NULL
);

CREATE TABLE public.account (
    id character(36) NOT NULL,
    branch_id character(36) DEFAULT ''::bpchar NOT NULL,
    number character varying(200) NOT NULL,
    customer_id character(36) DEFAULT ''::bpchar NOT NULL,
    account_type character varying(28) NOT NULL,
    target double precision DEFAULT 0 NOT NULL,
    target_info character varying(200) DEFAULT ''::character varying NOT NULL,
    sales_rep_id character(36) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint,
    balance double precision DEFAULT 0 NOT NULL,
    last_payment_date bigint DEFAULT 0 NOT NULL
);

CREATE TABLE public.bank_deposit (
    id character(36) NOT NULL,
    bank_account_id character(36) NOT NULL,
    amount double precision NOT NULL,
    date bigint NOT NULL
);

CREATE TABLE public.branch (
    id character(36) NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint
);

CREATE TABLE public.brand (
    id character(36) NOT NULL,
    name character varying(256) NOT NULL,
    code character(4) NOT NULL,
    logo character varying(128) NOT NULL
);

CREATE TABLE public.category (
    id character(36) NOT NULL,
    name character varying(256) NOT NULL
);

CREATE TABLE public.customer (
    id character(36) NOT NULL,
    branch_id character(36) DEFAULT ''::bpchar NOT NULL,
    email character varying(200) NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL,
    phone_number character varying(200) NOT NULL,
    address character varying(256) NOT NULL,
    sales_rep_id character(36) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint
);

CREATE TABLE public.daily_summary (
    income double precision NOT NULL,
    expenditure double precision NOT NULL,
    bank_deposit double precision NOT NULL,
    date bigint NOT NULL
);

CREATE TABLE public.ds_commission (
    id character(36) NOT NULL,
    account_id character(36) DEFAULT ''::bpchar NOT NULL,
    customer_id character(36) NOT NULL,
    amount double precision NOT NULL,
    date bigint NOT NULL,
    effective_date bigint DEFAULT 0 NOT NULL
);

CREATE TABLE public.expenditure (
    id character(36) NOT NULL,
    amount double precision NOT NULL,
    date bigint NOT NULL,
    reason character varying(200) DEFAULT ''::character varying NOT NULL
);

CREATE TABLE public.inventory (
    id character(36) NOT NULL,
    product_id character(36) DEFAULT ''::bpchar NOT NULL,
    branch_id character(36) DEFAULT ''::bpchar NOT NULL,
    tx_type character varying(33) NOT NULL,
    opening_balance double precision NOT NULL,
    quantity double precision DEFAULT 0 NOT NULL,
    narration character varying(200) DEFAULT ''::character varying NOT NULL,
    sales_rep_id character(36) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint
);

CREATE TABLE public.payment (
    id character(36) NOT NULL,
    sale_id character(36) NOT NULL,
    amount double precision NOT NULL,
    payment_method public.payment_method NOT NULL,
    sales_rep_id character(36) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint
);

CREATE TABLE public.product (
    id character(36) NOT NULL,
    brand_id character(36),
    category_id character(36) NOT NULL,
    name character varying(256) NOT NULL,
    description character varying(512) NOT NULL,
    sku character varying(128) NOT NULL,
    barcode character varying(128) NOT NULL,
    price double precision NOT NULL,
    reorder_level integer NOT NULL,
    image character varying(128),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    archived_at timestamp with time zone,
    created_by_id character(36) NOT NULL,
    updated_by_id character(36) NOT NULL,
    archived_by_id character(36),
    stock_balance integer NOT NULL
);

CREATE TABLE public.product_category (
    id character(36) NOT NULL,
    product_id character(36) NOT NULL,
    category_id character(36) NOT NULL
);

CREATE TABLE public.profit (
    id character(36) NOT NULL,
    amount double precision NOT NULL,
    narration character varying(200) DEFAULT ''::character varying NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint
);

CREATE TABLE public.reps_expense (
    id character(36) NOT NULL,
    sales_rep_id character(36) NOT NULL,
    amount double precision NOT NULL,
    reason character varying(200) NOT NULL,
    date bigint NOT NULL
);

CREATE TABLE public.sale (
    id character(36) NOT NULL,
    branch_id character(36) DEFAULT ''::bpchar NOT NULL,
    receipt_number character varying(128) NOT NULL,
    amount double precision NOT NULL,
    amount_tender double precision NOT NULL,
    balance double precision NOT NULL,
    customer_name character varying(256),
    phone_number character varying(28),
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint,
    created_by_id character(36) NOT NULL,
    updated_by_id character(36),
    archived_by_id character(36)
);

CREATE TABLE public.sale_item (
    id character(36) NOT NULL,
    sale_id character(36) NOT NULL,
    product_id character(36) NOT NULL,
    quantity integer NOT NULL,
    unit_price double precision NOT NULL,
    unit_cost_price double precision NOT NULL,
    stock_ids character varying(512) NOT NULL
);

CREATE TABLE public.stock (
    id character(36) NOT NULL,
    branch_id character(36) DEFAULT ''::bpchar NOT NULL,
    batch_number character varying(128) NOT NULL,
    product_id character(36) NOT NULL,
    unit_cost_price double precision NOT NULL,
    quantity integer NOT NULL,
    deducted_quantity integer NOT NULL,
    manufacture_date timestamp without time zone,
    expiry_date timestamp without time zone,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint,
    created_by_id character(36) NOT NULL,
    updated_by_id character(36) NOT NULL,
    archived_by_id character(36)
);

CREATE TABLE public.transaction (
    id character(36) NOT NULL,
    account_id character(36) DEFAULT ''::bpchar NOT NULL,
    tx_type character varying(33) NOT NULL,
    opening_balance double precision NOT NULL,
    amount double precision DEFAULT 0 NOT NULL,
    narration character varying(200) DEFAULT ''::character varying NOT NULL,
    sales_rep_id character(36) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    archived_at bigint,
    receipt_no character varying(33) DEFAULT ''::character varying NOT NULL,
    effective_date bigint DEFAULT 0 NOT NULL,
    payment_method character(36) DEFAULT 'cash'::bpchar NOT NULL
);

CREATE TABLE public.users (
    id character(36) NOT NULL,
    branch_id character(36) DEFAULT ''::bpchar NOT NULL,
    email character varying(200) NOT NULL,
    first_name character varying(200) DEFAULT ''::character varying NOT NULL,
    password_hash character varying(256) NOT NULL,
    password_salt character varying(36) NOT NULL,
    password_reset character varying(36) DEFAULT NULL::character varying,
    timezone character varying(128),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    last_name character varying(200) DEFAULT ''::character varying NOT NULL,
    phone_number character(36) DEFAULT ''::bpchar NOT NULL
);
