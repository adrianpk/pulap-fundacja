CREATE TABLE samples
	 id VARCHAR(36) PRIMARY KEY UNIQUE,
	 name VARCHAR(64) UNIQUE,
	 description TEXT,
	 user_id VARCHAR(64),
	 annotations JSONB NULL,
	 geolocation GEOGRAPHY(Point,4326) NULL,
	 started_at TIMESTAMP WITH TIME ZONE,
	 created_by VARCHAR(64),
	 is_active BOOLEAN,
	 is_logical_deleted BOOLEAN,
	 created_at TIMESTAMP WITH TIME ZONE,
	 updated_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE Samples
	ADD CONSTRAINT users_id_fkey
	FOREIGN KEY (user_id)
	REFERENCES samples
	ON DELETE CASCADE;

