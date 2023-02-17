#!/bin/bash
set -o errexit
set -x

SOURCE_DB="highload_architect-db-1"
REPLICA_DB="highload_architect-db_replica-1"

mysql_exec () {
  db=${1}
  query=${2}
  docker exec "${db}" mysql --user=root --password=password -e "${query}";
}

while ! (mysql_exec ${SOURCE_DB} "SELECT 1") >/dev/null; do
  echo "Waiting for source database ${SOURCE_DB} connection..."
  sleep 5
done

while ! (mysql_exec ${REPLICA_DB} "SELECT 1") >/dev/null; do
  echo "Waiting for replica database ${REPLICA_DB} connection..."
  sleep 5
done

echo "creating replica user"
REPLICATION_USER="replica"
REPLICATION_PASSWORD="replica"
mysql_exec ${SOURCE_DB} \
  "CREATE USER '$REPLICATION_USER'@'%' IDENTIFIED WITH mysql_native_password BY '$REPLICATION_PASSWORD';"
mysql_exec ${SOURCE_DB} \
  "GRANT REPLICATION SLAVE ON *.* TO '$REPLICATION_USER'@'%';"

echo "getting Source File and Position"
SOURCE_FILE="$(mysql_exec ${SOURCE_DB} 'show master status \G' | grep File | sed -n -e 's/^.*: //p')"
SOURCE_POSITION="$(mysql_exec ${SOURCE_DB} 'show master status \G' | grep Position | grep -Eo '[0-9]{1,}')"

echo "init replica server for connecting to the master server"
mysql_exec ${REPLICA_DB} \
  "CHANGE REPLICATION SOURCE TO \
  SOURCE_HOST = '${SOURCE_DB}', \
  SOURCE_PORT = 3306, \
  SOURCE_USER = '${REPLICATION_USER}', \
  SOURCE_PASSWORD = '${REPLICATION_PASSWORD}', \
  SOURCE_LOG_FILE = '${SOURCE_FILE}', \
  SOURCE_LOG_POS = ${SOURCE_POSITION};"

mysql_exec ${REPLICA_DB} "show replica status \G"
mysql_exec ${REPLICA_DB} "start replica";