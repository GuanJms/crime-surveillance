CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE crime_status AS ENUM (
    'NEW',
    'ASSIGNED',
    'RESOLVED'
);

CREATE TYPE patrol_status AS ENUM (
    'AVAILABLE',
    'BUSY',
    'OFF_DUTY'
);

CREATE TYPE user_role AS ENUM (
  'CITIZEN',
  'PATROL',
  'DISPATCHER',
  'ADMIN'
);


CREATE TABLE users(
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	username TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	role user_role NOT NULL DEFAULT 'CITIZEN',
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
	last_login TIMESTAMP WITHOUT TIME ZONE,
	last_activity TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE patrol_profile (
user_id UUID PRIMARY KEY,
officer_id TEXT UNIQUE NOT NULL,
officer_name TEXT NOT NULL,
status patrol_status NOT NULL DEFAULT 'AVAILABLE',
street TEXT,
city TEXT,
state TEXT,
latitude DOUBLE PRECISION,
longitude DOUBLE PRECISION,
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE crime (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
reporter_id UUID NOT NULL,
patrol_id UUID,
description TEXT,
status crime_status NOT NULL DEFAULT 'NEW',
street TEXT,
city TEXT,
state TEXT,
latitude DOUBLE PRECISION,
longitude DOUBLE PRECISION,
reported_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
CONSTRAINT fk_user FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,
CONSTRAINT fk_patrol FOREIGN KEY (patrol_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE event_log(
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL,
	service TEXT, 
	message JSONB,
	created_at TIMESTAMP DEFAULT now(),
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);