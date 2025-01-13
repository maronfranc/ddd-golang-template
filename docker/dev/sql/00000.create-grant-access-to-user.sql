-- Script to run the following command with variables and catching exceptions.
-- CREATE ROLE _user LOGIN PASSWORD _password;
-- CREATE DATABASE _db;
-- GRANT ALL PRIVILEGES ON DATABASE _db TO _user;
DO $$
DECLARE
	_db TEXT := 'db_123';
	_user TEXT := 'user_123';
	_password TEXT := 'pass123';
BEGIN

-- SEE: https://stackoverflow.com/a/55954480
EXECUTE FORMAT($f$CREATE ROLE %I LOGIN password '%s'$f$, _user, _password);
EXCEPTION WHEN duplicate_object THEN RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;

-- SEE: https://stackoverflow.com/a/36218838
CREATE EXTENSION IF NOT EXISTS dblink; -- enable extension 
IF EXISTS (SELECT 1 FROM pg_database WHERE datname = _db) THEN
   RAISE NOTICE 'Database already exists';
ELSE
   PERFORM dblink_connect('host=localhost user=' || _user || ' password=' || _password || ' dbname=' || current_database());
   PERFORM dblink_exec('CREATE DATABASE ' || _db);
END IF;

EXECUTE FORMAT($f$GRANT ALL PRIVILEGES ON DATABASE %I TO %s $f$, _db, _user);

END
$$
