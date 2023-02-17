#!/bin/bash
set -o errexit
set -x

SOURCE_DB="highload_architect-db-1"
REPLICA_1_DB="highload_architect-db_replica_1-1"
REPLICA_2_DB="highload_architect-db_replica_2-1"

mysql_exec () {
  db_num=$(($# - 1))
  query=${@: -1}
  for db in ${*: 1: ${db_num}}; do
    docker exec "${db}" mysql --user=root --password=password -e "${query}";
  done
}

while ! (mysql_exec ${SOURCE_DB} "SELECT 1") >/dev/null; do
  echo "Waiting for source database ${SOURCE_DB} connection..."
  sleep 5
done

while ! (mysql_exec ${REPLICA_1_DB} "SELECT 1") >/dev/null; do
  echo "Waiting for replica database ${REPLICA_1_DB} connection..."
  sleep 5
done

while ! (mysql_exec ${REPLICA_2_DB} "SELECT 1") >/dev/null; do
  echo "Waiting for replica database ${REPLICA_2_DB} connection..."
  sleep 5
done

echo "creating replica user"
REPLICATION_USER="replica"
REPLICATION_PASSWORD="replica"
mysql_exec ${SOURCE_DB} \
  "CREATE USER '$REPLICATION_USER'@'%' IDENTIFIED WITH mysql_native_password BY '$REPLICATION_PASSWORD';"
mysql_exec ${SOURCE_DB} \
  "GRANT REPLICATION SLAVE ON *.* TO '$REPLICATION_USER'@'%';"

echo "init replica server for connecting to the master server"
mysql_exec ${REPLICA_1_DB} ${REPLICA_2_DB} \
  "CHANGE REPLICATION SOURCE TO \
  SOURCE_HOST = '${SOURCE_DB}', \
  SOURCE_PORT = 3306, \
  SOURCE_USER = '${REPLICATION_USER}', \
  SOURCE_PASSWORD = '${REPLICATION_PASSWORD}', \
  SOURCE_AUTO_POSITION = 1;"

mysql_exec ${REPLICA_1_DB} ${REPLICA_2_DB} "show replica status \G"
mysql_exec ${REPLICA_1_DB} ${REPLICA_2_DB} "start replica";

echo "installing semisync plugin"
mysql_exec ${SOURCE_DB} \
  "INSTALL PLUGIN rpl_semi_sync_source SONAME 'semisync_source.so'; \
  INSTALL PLUGIN rpl_semi_sync_replica SONAME 'semisync_replica.so'; \
  SET GLOBAL rpl_semi_sync_source_enabled = 1;"

mysql_exec ${REPLICA_1_DB} ${REPLICA_2_DB} \
  "INSTALL PLUGIN rpl_semi_sync_source SONAME 'semisync_source.so'; \
  INSTALL PLUGIN rpl_semi_sync_replica SONAME 'semisync_replica.so'; \
  SET GLOBAL rpl_semi_sync_replica_enabled = 1; \
  STOP REPLICA IO_THREAD; \
  START REPLICA IO_THREAD;"

mysql_exec ${SOURCE_DB} "SHOW VARIABLES LIKE 'rpl_semi_sync%'"
mysql_exec ${REPLICA_1_DB} ${REPLICA_2_DB} "SHOW VARIABLES LIKE 'rpl_semi_sync%'"