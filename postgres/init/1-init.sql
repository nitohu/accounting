-- Create user for accounting application
CREATE USER accounting WITH PASSWORD 'supersecret';
ALTER USER accounting CREATEDB SUPERUSER;

CREATE DATABASE accounting OWNER "accounting";

