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

CREATE TABLE user_roles
(id UUID PRIMARY KEY,
 name VARCHAR(128),
 description VARCHAR(255) NULL,
 organization_id UUID,
 user_id UUID,
 role_id UUID,
 created_by UUID NULL,
 is_active BOOLEAN,
 is_logical_deleted BOOLEAN,
 created_at TIMESTAMP WITH TIME ZONE,
 updated_at TIMESTAMP WITH TIME ZONE);

ALTER TABLE user_roles
 ADD CONSTRAINT organization_id_fkey
 FOREIGN KEY (organization_id)
 REFERENCES organizations
 ON DELETE CASCADE;

ALTER TABLE user_roles
  ADD CONSTRAINT user_id_fkey
  FOREIGN KEY (user_id)
  REFERENCES users
  ON DELETE CASCADE;

ALTER TABLE user_roles
  ADD CONSTRAINT role_id_fkey
  FOREIGN KEY (role_id)
  REFERENCES roles
  ON DELETE CASCADE;
