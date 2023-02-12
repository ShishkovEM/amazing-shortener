CREATE TABLE IF NOT EXISTS "urls" (
    "short_uri" TEXT COLLATE pg_catalog.default NOT NULL,
    "original_url" TEXT COLLATE pg_catalog.default NOT NULL,
    "user_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "urls_pkey" PRIMARY KEY ("original_url")
);