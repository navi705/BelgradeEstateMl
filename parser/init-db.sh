#!/bin/bash
set -e

# This script creates a second user and database using environment variables
# provided via docker-compose/env file.

# 1. Create user and grant DB permissions (connecting to system 'postgres' db)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = '$PROJECT_USER') THEN
            CREATE USER $PROJECT_USER WITH PASSWORD '$PROJECT_PASSWORD';
        END IF;
    END
    \$\$;
    GRANT ALL PRIVILEGES ON DATABASE "$POSTGRES_DB" TO "$PROJECT_USER";
EOSQL

# 2. Grant schema permissions (connecting directly to our app database)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Grant schema permissions
    GRANT USAGE, CREATE ON SCHEMA public TO "$PROJECT_USER";
    ALTER SCHEMA public OWNER TO "$PROJECT_USER";
    
    -- Ensure future permissions
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "$PROJECT_USER";
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "$PROJECT_USER";
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "$PROJECT_USER";
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "$PROJECT_USER";
EOSQL
