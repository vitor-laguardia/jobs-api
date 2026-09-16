BEGIN;

CREATE TABLE users (
	id UUID DEFAULT gen_random_uuid(),
	name VARCHAR(256) NOT NULL,
	email VARCHAR(256) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY(id)
);

CREATE TABLE jobs (
	id UUID DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL,
	title VARCHAR(100) NOT NULL,
	description VARCHAR(256),
	status VARCHAR(20) NOT NULL CHECK (status IN ('pending','running', 'done', 'failed')),
	priority SMALLINT NOT NULL CHECK (priority IN (1 , 2, 3)),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY(id),
	CONSTRAINT fk_jobs_users
		FOREIGN KEY(user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
);

INSERT INTO users (name, email) VALUES ('alfred', 'alfred@gmail.com');
INSERT INTO users (name, email) VALUES ('derek', 'derek@gmail.com');

COMMIT;
