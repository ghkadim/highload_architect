#!/bin/bash
set -o errexit
set -x

SOURCE_DB="highload_architect-db_replica_1-1"
REPLICA_DB="highload_architect-db_replica_2-1"

REPLICATION_USER="replica"
REPLICATION_PASSWORD="replica"

mysql_exec () {
  db_num=$(($# - 1))
  query=${@: -1}
  for db in ${*: 1: ${db_num}}; do
    docker exec "${db}" mysql --user=root --password=password -e "${query}";
  done
}

echo "Stop io threads"
mysql_exec ${SOURCE_DB} ${REPLICA_DB} "STOP REPLICA IO_THREAD"
mysql_exec ${SOURCE_DB} ${REPLICA_DB} "SHOW PROCESSLIST "

echo "Stop replication on new source"
mysql_exec ${SOURCE_DB} "STOP REPLICA"
mysql_exec ${SOURCE_DB} "RESET MASTER"
mysql_exec ${SOURCE_DB} \
  "SET GLOBAL rpl_semi_sync_replica_enabled = 0;
  SET GLOBAL rpl_semi_sync_source_enabled = 1;"

echo "init replica server for connecting to the master server"
mysql_exec ${REPLICA_DB} "STOP REPLICA"
mysql_exec ${REPLICA_DB} \
  "CHANGE REPLICATION SOURCE TO \
  SOURCE_HOST = '${SOURCE_DB}', \
  SOURCE_PORT = 3306, \
  SOURCE_USER = '${REPLICATION_USER}', \
  SOURCE_PASSWORD = '${REPLICATION_PASSWORD}', \
  SOURCE_AUTO_POSITION = 1;"

mysql_exec ${REPLICA_DB} "SHOW REPLICA STATUS \G"
mysql_exec ${REPLICA_DB} "START REPLICA";

mysql_exec ${SOURCE_DB} "SHOW VARIABLES LIKE 'rpl_semi_sync%'"
mysql_exec ${REPLICA_DB} "SHOW VARIABLES LIKE 'rpl_semi_sync%'"


