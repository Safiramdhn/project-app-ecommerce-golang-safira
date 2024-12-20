--
-- PostgreSQL database dump
--

-- Dumped from database version 17.0
-- Dumped by pg_dump version 17.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: cart_status_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.cart_status_enum AS ENUM (
    'active',
    'checkout'
);


ALTER TYPE public.cart_status_enum OWNER TO postgres;

--
-- Name: order_status_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.order_status_enum AS ENUM (
    'on_progress',
    'success',
    'failed'
);


ALTER TYPE public.order_status_enum OWNER TO postgres;

--
-- Name: status_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.status_enum AS ENUM (
    'active',
    'deleted'
);


ALTER TYPE public.status_enum OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.addresses (
    id integer NOT NULL,
    user_id character varying NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    name character varying(100) NOT NULL,
    street character varying(255) NOT NULL,
    district character varying(100),
    city character varying(100),
    state character varying(100),
    postal_code character varying(20),
    country character varying(100) NOT NULL,
    status public.status_enum DEFAULT 'active'::public.status_enum NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.addresses OWNER TO postgres;

--
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.addresses_id_seq OWNER TO postgres;

--
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- Name: cart_item_variants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cart_item_variants (
    id integer NOT NULL,
    cart_item_id integer,
    item_variant_id integer,
    option_id integer,
    additional_price numeric(10,2),
    status public.status_enum DEFAULT 'active'::public.status_enum NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.cart_item_variants OWNER TO postgres;

--
-- Name: cart_item_variants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cart_item_variants_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cart_item_variants_id_seq OWNER TO postgres;

--
-- Name: cart_item_variants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cart_item_variants_id_seq OWNED BY public.cart_item_variants.id;


--
-- Name: cart_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cart_items (
    id integer NOT NULL,
    cart_id integer,
    product_id integer,
    amount integer DEFAULT 1,
    sub_total numeric(10,2) DEFAULT 0,
    status public.status_enum DEFAULT 'active'::public.status_enum NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.cart_items OWNER TO postgres;

--
-- Name: cart_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cart_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cart_items_id_seq OWNER TO postgres;

--
-- Name: cart_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cart_items_id_seq OWNED BY public.cart_items.id;


--
-- Name: carts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.carts (
    id integer NOT NULL,
    user_id character varying,
    total_amount integer DEFAULT 0,
    total_price numeric(10,2) DEFAULT 0,
    cart_status public.cart_status_enum DEFAULT 'active'::public.cart_status_enum,
    status public.status_enum DEFAULT 'active'::public.status_enum NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.carts OWNER TO postgres;

--
-- Name: carts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.carts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.carts_id_seq OWNER TO postgres;

--
-- Name: carts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.carts_id_seq OWNED BY public.carts.id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying NOT NULL,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: category_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.category_id_seq OWNER TO postgres;

--
-- Name: order_item_variants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_item_variants (
    id integer NOT NULL,
    order_item_id integer,
    variant_id integer,
    option_id integer,
    additional_price numeric(10,2),
    status public.status_enum DEFAULT 'active'::public.status_enum NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.order_item_variants OWNER TO postgres;

--
-- Name: order_item_variants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_item_variants_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_item_variants_id_seq OWNER TO postgres;

--
-- Name: order_item_variants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_item_variants_id_seq OWNED BY public.order_item_variants.id;


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_items (
    id integer NOT NULL,
    order_id integer,
    product_id integer,
    amount integer DEFAULT 0,
    subtotal numeric(10,2) DEFAULT 0,
    status public.status_enum DEFAULT 'active'::public.status_enum NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.order_items OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    user_id character varying,
    address_id integer,
    shipping_type text NOT NULL,
    shipping_cost numeric(10,2) DEFAULT 0,
    payment_method text NOT NULL,
    total_amount integer DEFAULT 0,
    total_price numeric(10,2) DEFAULT 0,
    order_status public.order_status_enum DEFAULT 'on_progress'::public.order_status_enum,
    status public.status_enum DEFAULT 'active'::public.status_enum NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id integer NOT NULL,
    name character varying NOT NULL,
    description character varying NOT NULL,
    category_id integer,
    price numeric(10,2) NOT NULL,
    discount numeric(10,2) DEFAULT 0,
    rating numeric(2,1) DEFAULT 0,
    photo_url text,
    has_variant boolean DEFAULT false,
    total_stock integer DEFAULT 0,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    CONSTRAINT products_rating_check CHECK (((rating >= (1)::numeric) AND (rating <= (5)::numeric)))
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: recommendations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.recommendations (
    id integer NOT NULL,
    product_id integer,
    is_recommended boolean DEFAULT false,
    set_in_banner boolean DEFAULT false,
    title character varying NOT NULL,
    subtitle character varying NOT NULL,
    photo_url text,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.recommendations OWNER TO postgres;

--
-- Name: recommendations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.recommendations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recommendations_id_seq OWNER TO postgres;

--
-- Name: recommendations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.recommendations_id_seq OWNED BY public.recommendations.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id character varying NOT NULL,
    name character varying NOT NULL,
    email character varying,
    phone_number character varying,
    password character varying NOT NULL,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: variation_options; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.variation_options (
    id integer NOT NULL,
    variation_id integer,
    option_value character varying NOT NULL,
    additional_price numeric(10,2) DEFAULT 0,
    stock integer DEFAULT 0,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.variation_options OWNER TO postgres;

--
-- Name: variation_options_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.variation_options_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.variation_options_id_seq OWNER TO postgres;

--
-- Name: variation_options_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.variation_options_id_seq OWNED BY public.variation_options.id;


--
-- Name: variations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.variations (
    id integer NOT NULL,
    product_id integer,
    attribute_name character varying NOT NULL,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.variations OWNER TO postgres;

--
-- Name: variations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.variations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.variations_id_seq OWNER TO postgres;

--
-- Name: variations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.variations_id_seq OWNED BY public.variations.id;


--
-- Name: weekly_promos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.weekly_promos (
    id integer NOT NULL,
    product_id integer,
    promo_discount numeric(10,2) DEFAULT 0 NOT NULL,
    start_date date,
    end_date date,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.weekly_promos OWNER TO postgres;

--
-- Name: weekly_promos_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.weekly_promos_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.weekly_promos_id_seq OWNER TO postgres;

--
-- Name: weekly_promos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.weekly_promos_id_seq OWNED BY public.weekly_promos.id;


--
-- Name: wishlist; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wishlist (
    id integer NOT NULL,
    user_id character varying,
    product_id integer,
    status public.status_enum DEFAULT 'active'::public.status_enum,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


ALTER TABLE public.wishlist OWNER TO postgres;

--
-- Name: wishlist_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.wishlist_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.wishlist_id_seq OWNER TO postgres;

--
-- Name: wishlist_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.wishlist_id_seq OWNED BY public.wishlist.id;


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Name: cart_item_variants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_item_variants ALTER COLUMN id SET DEFAULT nextval('public.cart_item_variants_id_seq'::regclass);


--
-- Name: cart_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items ALTER COLUMN id SET DEFAULT nextval('public.cart_items_id_seq'::regclass);


--
-- Name: carts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts ALTER COLUMN id SET DEFAULT nextval('public.carts_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: order_item_variants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_variants ALTER COLUMN id SET DEFAULT nextval('public.order_item_variants_id_seq'::regclass);


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: recommendations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recommendations ALTER COLUMN id SET DEFAULT nextval('public.recommendations_id_seq'::regclass);


--
-- Name: variation_options id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variation_options ALTER COLUMN id SET DEFAULT nextval('public.variation_options_id_seq'::regclass);


--
-- Name: variations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variations ALTER COLUMN id SET DEFAULT nextval('public.variations_id_seq'::regclass);


--
-- Name: weekly_promos id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weekly_promos ALTER COLUMN id SET DEFAULT nextval('public.weekly_promos_id_seq'::regclass);


--
-- Name: wishlist id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlist ALTER COLUMN id SET DEFAULT nextval('public.wishlist_id_seq'::regclass);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.addresses (id, user_id, is_default, name, street, district, city, state, postal_code, country, status, created_at, updated_at, deleted_at) FROM stdin;
1	c63b1125-a7d8-4ca0-b911-a866167345b0	f	House 2	New street	district A	New city	New state	12345	New country	active	2024-11-23 11:38:56.955833	2024-11-24 20:08:05.67485	\N
2	c63b1125-a7d8-4ca0-b911-a866167345b0	t	Office	New street	New district	New city	New state	12345	New country	active	2024-11-23 11:39:36.220399	2024-11-24 20:08:05.684449	\N
3	c63b1125-a7d8-4ca0-b911-a866167345b0	f	House 2	New street	\N	\N	\N	12345	New country	deleted	2024-11-23 12:49:30.855129	2024-11-23 12:49:30.855129	2024-11-24 20:09:03.706894
\.


--
-- Data for Name: cart_item_variants; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cart_item_variants (id, cart_item_id, item_variant_id, option_id, additional_price, status, created_at, updated_at, deleted_at) FROM stdin;
1	2	1	1	0.00	active	2024-11-24 21:01:55.534027	\N	\N
2	2	2	5	0.00	active	2024-11-24 21:01:55.543338	\N	\N
\.


--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cart_items (id, cart_id, product_id, amount, sub_total, status, created_at, updated_at, deleted_at) FROM stdin;
2	1	1	1	15.20	deleted	2024-11-24 21:01:55.511249	\N	2024-11-24 21:26:07.662455
3	1	3	2	25.60	active	2024-11-24 21:03:21.542341	2024-11-24 21:58:13.566668	\N
\.


--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.carts (id, user_id, total_amount, total_price, cart_status, status, created_at, updated_at, deleted_at) FROM stdin;
1	c63b1125-a7d8-4ca0-b911-a866167345b0	2	25.60	checkout	active	2024-11-24 20:52:21.152505	2024-11-24 21:58:13.603331	\N
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, status, created_at, updated_at, deleted_at) FROM stdin;
1	Electronics	active	2024-11-20 20:56:09.563013	\N	\N
2	Books	active	2024-11-20 20:56:09.563013	\N	\N
3	Clothing	active	2024-11-20 20:56:09.563013	\N	\N
4	Home Appliances	active	2024-11-20 20:56:09.563013	\N	\N
5	Toys	active	2024-11-20 20:56:09.563013	\N	\N
6	Groceries	active	2024-11-20 20:56:09.563013	\N	\N
7	Sports	active	2024-11-20 20:56:09.563013	\N	\N
8	Health & Beauty	active	2024-11-20 20:56:09.563013	\N	\N
9	Automotive	active	2024-11-20 20:56:09.563013	\N	\N
10	Furniture	active	2024-11-20 20:56:09.563013	\N	\N
11	Jewelry	active	2024-11-20 20:56:09.563013	\N	\N
12	Music	active	2024-11-20 20:56:09.563013	\N	\N
13	Stationery	active	2024-11-20 20:56:09.563013	\N	\N
14	Gardening	active	2024-11-20 20:56:09.563013	\N	\N
15	Pet Supplies	active	2024-11-20 20:56:09.563013	\N	\N
\.


--
-- Data for Name: order_item_variants; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_item_variants (id, order_item_id, variant_id, option_id, additional_price, status, created_at) FROM stdin;
1	2	1	1	0.00	active	2024-11-24 22:30:36.544083
2	2	2	5	0.00	active	2024-11-24 22:30:36.554379
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_items (id, order_id, product_id, amount, subtotal, status, created_at) FROM stdin;
1	1	1	1	15.20	active	2024-11-24 22:28:43.661296
2	2	1	1	15.20	active	2024-11-24 22:30:36.541588
3	2	3	2	25.60	active	2024-11-24 22:30:36.555712
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, address_id, shipping_type, shipping_cost, payment_method, total_amount, total_price, order_status, status, created_at) FROM stdin;
1	c63b1125-a7d8-4ca0-b911-a866167345b0	2	regular	0.00	bank	2	25.60	failed	active	2024-11-24 22:28:43.569753
2	c63b1125-a7d8-4ca0-b911-a866167345b0	2	regular	0.00	bank	2	25.60	success	active	2024-11-24 22:30:36.537275
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, description, category_id, price, discount, rating, photo_url, has_variant, total_stock, status, created_at, updated_at, deleted_at) FROM stdin;
1	Casual T-Shirt	A comfortable cotton T-shirt, perfect for everyday wear.	3	15.99	2.00	4.5	https://example.com/tshirt.jpg	t	120	active	2024-11-13 15:41:48.315155	\N	\N
2	Formal Shirt	A sleek and stylish formal shirt, ideal for office and events.	3	25.99	5.00	4.8	https://example.com/shirt.jpg	t	33	active	2024-11-03 15:41:48.315155	\N	\N
3	Stuffed Bear	A cute and cuddly stuffed bear, perfect as a gift or decoration.	5	12.99	0.00	4.2	https://example.com/bear.jpg	f	100	active	2024-11-18 15:41:48.315155	\N	\N
4	Skin Care Kit	A premium skin care kit to rejuvenate and nourish your skin.	8	45.99	10.00	4.7	https://example.com/skincare.jpg	t	45	active	2024-11-08 15:41:48.315155	\N	\N
5	Wooden Chair	A sturdy wooden chair, perfect for dining or working.	10	89.99	0.00	4.3	https://example.com/chair.jpg	f	100	active	2024-10-14 15:41:48.315155	\N	\N
\.


--
-- Data for Name: recommendations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.recommendations (id, product_id, is_recommended, set_in_banner, title, subtitle, photo_url, status, created_at, updated_at, deleted_at) FROM stdin;
1	1	t	f	Casual Comfort	Perfect for daily wear and casual outings.	https://example.com/tshirt.jpg	active	2024-11-21 19:59:21.817842	\N	\N
2	2	t	t	Formal Elegance	Upgrade your wardrobe with sleek style.	https://example.com/shirt.jpg	active	2024-11-21 19:59:21.817842	\N	\N
3	3	t	f	Gift of Cuteness	An adorable gift for your loved ones.	https://example.com/bear.jpg	active	2024-11-21 19:59:21.817842	\N	\N
4	4	t	t	Skin Care Deluxe	Rejuvenate your skin with premium care.	https://example.com/skincare.jpg	active	2024-11-21 19:59:21.817842	\N	\N
5	5	f	f	Elegant Wood Design	Perfect for your dining or working space.	https://example.com/chair.jpg	active	2024-11-21 19:59:21.817842	\N	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, phone_number, password, status, created_at, updated_at, deleted_at) FROM stdin;
c63b1125-a7d8-4ca0-b911-a866167345b0	User Test	test@example.com	\N	$2a$10$daB21EZbTN4BHaLmS2ApNefi.3S2IJ64OOR/.s.Mg6xrergJC7hKG	active	2024-11-20 10:39:23.846118	\N	\N
c847e893-2efc-4032-8888-2922c461578e	User Test 2	test2@example.com	\N	$2a$10$QNGw/N53pmCoycyZUGkuseX5zfFDfhfSh6qG5J9iMThCR.YbbRKAq	active	2024-11-20 10:50:55.945414	\N	\N
5c3f63b7-970e-4b52-b643-2c922c812f58	User Test 3	\N	0987654321	$2a$10$AilIH7IMDAopkGTHPF3IqOcsBbY73aEwg1pND.aTKND3mZcxQRKTm	active	2024-11-20 10:59:12.981641	\N	\N
c8cfe426-4cf4-42a3-9744-ada5c3b71fa0	User Test 4	\N	1234567890	$2a$10$mX.lfEZBB1Vt8ZlCSbTU/.gasPS5JEOSvqZJdazOA4WnuerKy20Xe	active	2024-11-20 10:59:29.089459	\N	\N
e5d6e296-6803-4fbb-98c8-62d9d51f17a0	User Test 5	example@email.com	\N	$2a$10$QUpHcFeS6/XKLcDWcJk3eePrDIXx1V38JXmyjqo3f01QjlUBIO9ZG	active	2024-11-25 19:37:27.61118	\N	\N
\.


--
-- Data for Name: variation_options; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.variation_options (id, variation_id, option_value, additional_price, stock, status, created_at) FROM stdin;
1	1	Small	0.00	50	active	2024-11-23 15:42:02.295823
2	1	Medium	1.00	40	active	2024-11-23 15:42:02.295823
3	1	Large	2.00	30	active	2024-11-23 15:42:02.295823
4	2	Red	0.00	60	active	2024-11-23 15:42:02.295823
5	2	Blue	0.00	60	active	2024-11-23 15:42:02.295823
6	3	Small	0.00	15	active	2024-11-23 15:42:05.095837
7	3	Medium	1.50	10	active	2024-11-23 15:42:05.095837
8	3	Large	2.50	8	active	2024-11-23 15:42:05.095837
9	4	Lavender	2.00	25	active	2024-11-23 15:42:07.675364
10	4	Rose	2.50	20	active	2024-11-23 15:42:07.675364
\.


--
-- Data for Name: variations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.variations (id, product_id, attribute_name, status, created_at) FROM stdin;
1	1	Size	active	2024-11-23 15:41:57.324144
2	1	Color	active	2024-11-23 15:41:57.324144
3	2	Size	active	2024-11-23 15:41:57.324144
4	4	Fragrance	active	2024-11-23 15:41:57.324144
\.


--
-- Data for Name: weekly_promos; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.weekly_promos (id, product_id, promo_discount, start_date, end_date, status, created_at, updated_at, deleted_at) FROM stdin;
1	1	3.00	2024-11-17	2024-11-24	active	2024-11-22 20:32:33.836441	\N	\N
2	2	7.00	2024-11-12	2024-11-19	active	2024-11-22 20:32:33.836441	\N	\N
3	3	1.50	2024-11-20	2024-11-27	active	2024-11-22 20:32:33.836441	\N	\N
4	4	12.00	2024-11-15	2024-11-29	active	2024-11-22 20:32:33.836441	\N	\N
5	5	10.00	2024-11-07	2024-11-14	active	2024-11-22 20:32:33.836441	\N	\N
\.


--
-- Data for Name: wishlist; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wishlist (id, user_id, product_id, status, created_at, updated_at, deleted_at) FROM stdin;
1	c63b1125-a7d8-4ca0-b911-a866167345b0	1	active	2024-11-22 16:15:38.136386	\N	\N
2	c63b1125-a7d8-4ca0-b911-a866167345b0	2	deleted	2024-11-22 16:22:20.776564	\N	2024-11-22 16:32:03.259437
\.


--
-- Name: addresses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.addresses_id_seq', 3, true);


--
-- Name: cart_item_variants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.cart_item_variants_id_seq', 2, true);


--
-- Name: cart_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.cart_items_id_seq', 3, true);


--
-- Name: carts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.carts_id_seq', 1, true);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 15, true);


--
-- Name: category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.category_id_seq', 15, true);


--
-- Name: order_item_variants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_item_variants_id_seq', 2, true);


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_items_id_seq', 3, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 2, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 5, true);


--
-- Name: recommendations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.recommendations_id_seq', 5, true);


--
-- Name: variation_options_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.variation_options_id_seq', 10, true);


--
-- Name: variations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.variations_id_seq', 4, true);


--
-- Name: weekly_promos_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.weekly_promos_id_seq', 5, true);


--
-- Name: wishlist_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.wishlist_id_seq', 2, true);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: cart_item_variants cart_item_variants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_item_variants
    ADD CONSTRAINT cart_item_variants_pkey PRIMARY KEY (id);


--
-- Name: cart_items cart_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_pkey PRIMARY KEY (id);


--
-- Name: carts carts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_pkey PRIMARY KEY (id);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: order_item_variants order_item_variants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_variants
    ADD CONSTRAINT order_item_variants_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: recommendations recommendations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recommendations
    ADD CONSTRAINT recommendations_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: variation_options variation_options_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variation_options
    ADD CONSTRAINT variation_options_pkey PRIMARY KEY (id);


--
-- Name: variations variations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variations
    ADD CONSTRAINT variations_pkey PRIMARY KEY (id);


--
-- Name: weekly_promos weekly_promos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weekly_promos
    ADD CONSTRAINT weekly_promos_pkey PRIMARY KEY (id);


--
-- Name: wishlist wishlist_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlist
    ADD CONSTRAINT wishlist_pkey PRIMARY KEY (id);


--
-- Name: addresses addresses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: cart_item_variants cart_item_variants_cart_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_item_variants
    ADD CONSTRAINT cart_item_variants_cart_item_id_fkey FOREIGN KEY (cart_item_id) REFERENCES public.cart_items(id) ON DELETE CASCADE;


--
-- Name: cart_item_variants cart_item_variants_item_variant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_item_variants
    ADD CONSTRAINT cart_item_variants_item_variant_id_fkey FOREIGN KEY (item_variant_id) REFERENCES public.variations(id) ON DELETE SET NULL;


--
-- Name: cart_item_variants cart_item_variants_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_item_variants
    ADD CONSTRAINT cart_item_variants_option_id_fkey FOREIGN KEY (option_id) REFERENCES public.variation_options(id) ON DELETE SET NULL;


--
-- Name: cart_items cart_items_cart_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES public.carts(id) ON DELETE CASCADE;


--
-- Name: cart_items cart_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE SET NULL;


--
-- Name: carts carts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: order_item_variants order_item_variants_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_variants
    ADD CONSTRAINT order_item_variants_option_id_fkey FOREIGN KEY (option_id) REFERENCES public.variation_options(id) ON DELETE SET NULL;


--
-- Name: order_item_variants order_item_variants_order_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_variants
    ADD CONSTRAINT order_item_variants_order_item_id_fkey FOREIGN KEY (order_item_id) REFERENCES public.order_items(id) ON DELETE CASCADE;


--
-- Name: order_item_variants order_item_variants_variant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_variants
    ADD CONSTRAINT order_item_variants_variant_id_fkey FOREIGN KEY (variant_id) REFERENCES public.variations(id) ON DELETE SET NULL;


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: orders orders_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE SET NULL;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- Name: variation_options variation_options_variation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variation_options
    ADD CONSTRAINT variation_options_variation_id_fkey FOREIGN KEY (variation_id) REFERENCES public.variations(id) ON DELETE CASCADE;


--
-- Name: variations variations_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variations
    ADD CONSTRAINT variations_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: wishlist wishlist_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlist
    ADD CONSTRAINT wishlist_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

