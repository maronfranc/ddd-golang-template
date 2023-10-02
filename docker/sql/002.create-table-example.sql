START TRANSACTION;

INSERT INTO migrations (file) VALUES ('002.create-table-example.sql');

CREATE TABLE IF NOT EXISTS examples (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
 	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
 	deleted_at TIMESTAMPTZ,
	title VARCHAR NOT NULL,
	description VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS example_sub_refs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
 	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
 	deleted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name VARCHAR NOT NULL,
	example_id UUID NOT NULL,
	FOREIGN KEY (example_id) REFERENCES examples(id)
);

COMMIT;