#!/bin/bash
set -o errexit
set -x

REPLICA_1_DB="highload_architect-db_1-1"
REPLICA_2_DB="highload_architect-db_2-1"

REPLICAS=(
  "$REPLICA_1_DB"
  "$REPLICA_2_DB"
)

REPLICA_SOURCES=(
  "$REPLICA_1_DB:$REPLICA_2_DB"
  "$REPLICA_2_DB:$REPLICA_1_DB"
)

mysql_exec () {
  db_num=$(($# - 1))
  query=${@: -1}
  for db in ${*: 1: ${db_num}}; do
    docker exec "${db}" mysql --user=root --password=password -e "${query}";
  done
}

for replica in "${REPLICAS[@]}"; do
  while ! (mysql_exec "${replica}" "SELECT 1") >/dev/null; do
    echo "Waiting for replica database ${replica} connection..."
    sleep 5
  done
done

echo "creating replica user"
REPLICATION_USER="replica"
REPLICATION_PASSWORD="replica"

for replica in "${REPLICAS[@]}"; do
  mysql_exec "${replica}" \
    "CREATE USER '$REPLICATION_USER'@'%' IDENTIFIED WITH mysql_native_password BY '$REPLICATION_PASSWORD';"
  mysql_exec "${replica}" \
    "GRANT REPLICATION SLAVE ON *.* TO '$REPLICATION_USER'@'%';"
  mysql_exec "${replica}" "CREATE USER 'haproxy_check'@'%';"
  mysql_exec "${replica}" "FLUSH PRIVILEGES;"
done

for replica_src in "${REPLICA_SOURCES[@]}"; do
  REPLICA_DB="${replica_src%%:*}"
  SOURCE_DB="${replica_src##*:}"
  echo "replica: ${REPLICA_DB}, source: ${SOURCE_DB}"
  echo "getting Source File and Position"
  SOURCE_FILE="$(mysql_exec "${SOURCE_DB}" 'show master status \G' | grep File | sed -n -e 's/^.*: //p')"
  SOURCE_POSITION="$(mysql_exec "${SOURCE_DB}" 'show master status \G' | grep Position | grep -Eo '[0-9]{1,}')"

  echo "init replica server for connecting to the master server"
  mysql_exec "${REPLICA_DB}" \
    "CHANGE REPLICATION SOURCE TO \
    SOURCE_HOST = '${SOURCE_DB}', \
    SOURCE_PORT = 3306, \
    SOURCE_USER = '${REPLICATION_USER}', \
    SOURCE_PASSWORD = '${REPLICATION_PASSWORD}', \
    SOURCE_LOG_FILE = '${SOURCE_FILE}', \
    SOURCE_LOG_POS = ${SOURCE_POSITION};"
done

for replica in "${REPLICAS[@]}"; do
  mysql_exec "${replica}" "show replica status \G"
  mysql_exec "${replica}" "start replica";
done