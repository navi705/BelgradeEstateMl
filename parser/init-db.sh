#!/bin/bash
set -e

# This script creates a second user and database using environment variables
# provided via docker-compose/env file.

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    CREATE USER $PROJECT_USER WITH PASSWORD '$PROJECT_PASSWORD';
    GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $PROJECT_USER;
    
    -- Connect to the local database to grant schema permissions
    \c $DB_NAME
    
    GRANT ALL ON SCHEMA public TO $PROJECT_USER;
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO $PROJECT_USER;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO $PROJECT_USER;
EOSQL
