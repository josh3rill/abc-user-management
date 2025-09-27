# #!/bin/bash
# set -euo pipefail

# echo "Starting ABC User Management Backend..."


# # Start the Go application in the background
# echo "Starting the application..."
# exec /app/main &
# APP_PID=$!

# # Wait for MySQL to be ready
# echo "Waiting for MySQL to be ready..."
# MYSQL_READY=0
# for i in {1..30}; do
#     if mysqladmin ping -h"${DB_HOST:-mysql}" -P"${DB_PORT:-3306}" -u"${DB_USER:-abc_user}" -p"${DB_PASSWORD:-v12345}" --silent; then
#         MYSQL_READY=1
#         break
#     fi
#     echo "MySQL not ready yet, retrying in 2s... ($i/30)"
#     sleep 2
# done
# if [ "$MYSQL_READY" -eq 1 ]; then
#     echo "MySQL is ready!"
# else
#     echo "Error: MySQL is not ready after waiting. Exiting."
#     kill $APP_PID
#     exit 1
# fi

# # Run migrations if present
# echo "Running database migrations..."
# if [ -d "./migrations" ]; then
#     for migration in ./migrations/*.sql; do
#         if [ -f "$migration" ]; then
#             echo "Applying migration: $(basename $migration)"
#             mysql -h"${DB_HOST:-mysql}" -P"${DB_PORT:-3306}" \
#                   -u"${DB_USER:-abc_user}" -p"${DB_PASSWORD:-v12345}" \
#                   "${DB_NAME:-abc_user_management}" < "$migration" || {
#                 echo "Warning: Migration $(basename $migration) may have already been applied or failed"
#             }
#         fi
#     done
#     echo "Migrations completed!"
# else
#     echo "No migrations directory found, skipping migrations"
# fi

# # Wait for the Go application to finish
# wait $APP_PID