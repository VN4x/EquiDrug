-- Diet macro tracking for travel-safe eating

CREATE TABLE macro_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    protein_g DOUBLE PRECISION NOT NULL DEFAULT 130,
    carbs_g DOUBLE PRECISION NOT NULL DEFAULT 80,
    fat_g DOUBLE PRECISION NOT NULL DEFAULT 25,
    calories_kcal DOUBLE PRECISION,
    notes TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE food_reference (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    name_local TEXT,
    locale TEXT DEFAULT 'en',
    country_code CHAR(2),
    serving_label TEXT NOT NULL DEFAULT '1 serving',
    protein_g DOUBLE PRECISION NOT NULL DEFAULT 0,
    carbs_g DOUBLE PRECISION NOT NULL DEFAULT 0,
    fat_g DOUBLE PRECISION NOT NULL DEFAULT 0,
    calories_kcal DOUBLE PRECISION,
    tags TEXT[],
    source TEXT NOT NULL DEFAULT 'curated',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE food_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    logged_on DATE NOT NULL DEFAULT CURRENT_DATE,
    food_name TEXT NOT NULL,
    serving_label TEXT,
    protein_g DOUBLE PRECISION NOT NULL DEFAULT 0,
    carbs_g DOUBLE PRECISION NOT NULL DEFAULT 0,
    fat_g DOUBLE PRECISION NOT NULL DEFAULT 0,
    calories_kcal DOUBLE PRECISION,
    confidence DOUBLE PRECISION NOT NULL DEFAULT 0.8,
    source TEXT NOT NULL DEFAULT 'manual',
    image_url TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_food_reference_name ON food_reference USING gin (to_tsvector('simple', name));
CREATE INDEX idx_food_reference_local ON food_reference USING gin (to_tsvector('simple', coalesce(name_local, '')));
CREATE INDEX idx_food_log_user_day ON food_log(user_id, logged_on);
