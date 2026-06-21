-- EquiDrug initial schema
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE product_type AS ENUM ('drug', 'supplement', 'vitamin', 'superfood', 'food');
CREATE TYPE item_category AS ENUM ('consume', 'avoid');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE,
    display_name TEXT,
    locale TEXT DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE active_ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inn TEXT NOT NULL,
    name TEXT NOT NULL,
    strength DOUBLE PRECISION NOT NULL DEFAULT 0,
    unit TEXT NOT NULL DEFAULT 'mg',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type product_type NOT NULL,
    category item_category NOT NULL DEFAULT 'consume',
    brand_name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    country_code CHAR(2) NOT NULL,
    image_url TEXT,
    retailer_hint TEXT,
    price_hint TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_ingredients (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES active_ingredients(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, ingredient_id)
);

CREATE TABLE equivalences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    origin_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    foreign_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    confidence DOUBLE PRECISION NOT NULL DEFAULT 0.8,
    notes TEXT,
    source TEXT NOT NULL DEFAULT 'curated',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (origin_product_id, foreign_product_id)
);

CREATE TABLE locker_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    custom_name TEXT NOT NULL,
    dose_per_day DOUBLE PRECISION NOT NULL DEFAULT 1,
    dose_unit TEXT NOT NULL DEFAULT 'tablet',
    frequency TEXT NOT NULL DEFAULT 'daily',
    category item_category NOT NULL DEFAULT 'consume',
    source_photo TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE trip_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    origin_country CHAR(2) NOT NULL,
    dest_country CHAR(2) NOT NULL,
    dest_city TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    spare_percent DOUBLE PRECISION NOT NULL DEFAULT 5,
    preferred_brands TEXT[],
    status TEXT NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE trip_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES trip_plans(id) ON DELETE CASCADE,
    locker_item_id UUID NOT NULL REFERENCES locker_items(id) ON DELETE CASCADE,
    origin_product_id UUID REFERENCES products(id),
    foreign_product_id UUID REFERENCES products(id),
    quantity_needed DOUBLE PRECISION NOT NULL,
    quantity_unit TEXT NOT NULL,
    estimated_price TEXT,
    where_to_buy TEXT,
    bought BOOLEAN NOT NULL DEFAULT false,
    notes TEXT
);

CREATE TABLE wiki_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    origin_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    foreign_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    confidence DOUBLE PRECISION NOT NULL DEFAULT 0.8,
    notes TEXT,
    source TEXT NOT NULL DEFAULT 'user_lookup',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, origin_product_id, foreign_product_id)
);

CREATE TABLE avoid_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    locale_from TEXT NOT NULL,
    term_from TEXT NOT NULL,
    locale_to TEXT NOT NULL,
    term_to TEXT NOT NULL,
    safe_alternate TEXT,
    severity TEXT NOT NULL DEFAULT 'warning',
    notes TEXT
);

CREATE INDEX idx_products_country ON products(country_code);
CREATE INDEX idx_products_display_name ON products USING gin (to_tsvector('simple', display_name));
CREATE INDEX idx_equivalences_origin ON equivalences(origin_product_id);
CREATE INDEX idx_locker_user ON locker_items(user_id);
CREATE INDEX idx_trip_user ON trip_plans(user_id);
