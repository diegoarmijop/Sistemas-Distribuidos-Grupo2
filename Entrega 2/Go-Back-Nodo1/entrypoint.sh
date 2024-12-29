#!/bin/sh

# Generar el archivo .env dinámicamente
cat <<EOF > .env
DB_HOST=${DB_HOST:-localhost}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-node_database}
DB_PORT=${DB_PORT:-5432}
DB_SSLMODE=${DB_SSLMODE:-disable}
DB_TIMEZONE=${DB_TIMEZONE:-UTC}
APP_PORT=${PORT:-8081}
RABBITMQ_URL=${RABBITMQ_URL:-amqp://guest:guest@rabbitmq:5672/}
BASE_CENTRAL_URL=${BASE_CENTRAL_URL:-http://localhost:8080}
EOF

# Mostrar el contenido del .env para verificarlo (opcional)
echo "Archivo .env generado:"
cat .env

# Ejecutar la aplicación
./main
