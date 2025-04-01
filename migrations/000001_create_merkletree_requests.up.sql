CREATE TABLE IF NOT EXISTS "merkletree_requests" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "network" VARCHAR(255) NOT NULL,
    "processed" BOOLEAN DEFAULT FALSE,
    "root" VARCHAR(255),
    "tree" JSONB
);
