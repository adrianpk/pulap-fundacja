-- Copyright (c) 2017 Kuguar <licenses@kuguar.io> Author: Adrian P.K. <apk@kuguar.io>
--
-- MIT License
--
-- Permission is hereby granted, free of charge, to any person obtaining
-- a copy of this software and associated documentation files (the
-- "Software"), to deal in the Software without restriction, including
-- without limitation the rights to use, copy, modify, merge, publish,
-- distribute, sublicense, and/or sell copies of the Software, and to
-- permit persons to whom the Software is furnished to do so, subject to
-- the following conditions:
--
-- The above copyright notice and this permission notice shall be
-- included in all copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
-- EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
-- MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
-- NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
-- LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
-- OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
-- WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

CREATE TABLE users
(id UUID PRIMARY KEY,
 username VARCHAR(32) UNIQUE,
 password_hash CHAR(128),
 email VARCHAR(255) UNIQUE,
 first_name VARCHAR(32),
 middle_names VARCHAR(32) NULL,
 last_name VARCHAR(64),
 started_at TIMESTAMP WITH TIME ZONE,
 card JSONB NULL,
 annotations JSONB NULL,
 geolocation GEOGRAPHY(Point,4326),
 created_by UUID NULL,
 is_active BOOLEAN,
 is_logical_deleted BOOLEAN,
 created_at TIMESTAMP WITH TIME ZONE,
 updated_at TIMESTAMP WITH TIME ZONE);


-- ALTER TABLE users
-- ADD COLUMN geolocation geography(Point,4326);
-- SELECT AddGeographyColumn ('public','accounts','location',4326,'POINT',2);
