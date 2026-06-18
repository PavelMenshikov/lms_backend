#!/bin/bash
set -euo pipefail

BACKUP_DIR="/root/backups"
RETENTION_DAYS=7
DB_CONTAINER="cap-edu-db-prod"
DB_USER="capuser"
DB_NAME="capedu"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/capedu_${TIMESTAMP}.sql.gz"

mkdir -p "$BACKUP_DIR"

echo "[$(date)] Starting DB backup..."

docker exec "$DB_CONTAINER" pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo "[$(date)] Backup created: $BACKUP_FILE ($BACKUP_SIZE)"
else
    echo "[$(date)] Backup FAILED!"
    rm -f "$BACKUP_FILE"
    exit 1
fi

echo "[$(date)] Cleaning backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "capedu_*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete

echo "[$(date)] Backup completed. Disk usage:"
df -h / | tail -1
