Create a database (First time)
-----------------------------------------------------

#### Install PostgreSQL on Mac

```bash
# Install
brew install postgresql@16

# Start
brew services start postgresql@16
```

#### Access PostgreSQL - User `vinh` (root user on Mac) 
```bash
psql -U vinh -d postgres

# Or Have hosting
psql -h localhost -U vinh -d postgres
```

#### List database
```postgresql
-- postgres:# 
SELECT datname FROM pg_database WHERE datistemplate = false;
```


#### Create database
```postgresql
-- postgres:#
CREATE DATABASE gfly;
\c gfly;
```

#### Create database
```postgresql
-- gfly-#
CREATE ROLE "user" WITH LOGIN PASSWORD 'secret';
ALTER ROLE "user" CREATEDB;
ALTER USER "user" WITH SUPERUSER;
GRANT ALL PRIVILEGES ON DATABASE gfly TO "user";
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA "public" TO "user";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA "public" TO "user";
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA "public" TO "user";
GRANT ALL ON SCHEMA public TO "user";
\q
```

#### Access PostgreSQL - User `user`
```bash
psql -U user -d gfly
```

#### Check tables in DB `gfly`
```postgresql
-- gfly-#
\dt
```

Re-create database
-----------------------------------------------------
#### Access PostgreSQL - User `user`
```bash
psql -h localhost -U user -d postgres
```

#### Reset current DB
```postgresql
-- postgres-#
REVOKE CONNECT ON DATABASE gfly FROM public;
SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'gfly';
DROP DATABASE IF EXISTS gfly;
CREATE DATABASE gfly;
\c gfly;
```
Continue below commands:
```postgresql
-- gfly:#
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA "public" TO "user";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA "public" TO "user";
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA "public" TO "user";
GRANT ALL ON SCHEMA public TO "user";
\q
```

Backup & Restore DB
-----------------------------------------------------
#### Backup

    pg_dump -h localhost -U user -d gfly_db > /Users/vinh/gfly_code/database/gfly_db.sql

#### Restore

    psql -U user -d gfly_db < /Users/vinh/gfly_code/database/gfly_db.sql

