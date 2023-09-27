#!/bin/bash
set -o errexit

REPLICA_1_DB="db_1"
REPLICA_2_DB="db_2"

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
    >&2 echo "${db}> ${query}"
    mysql -h "${db}" --user=root --password="${MYSQL_ROOT_PASSWORD:-password}" -e "${query}";
  done
}

for replica in "${REPLICAS[@]}"; do
  while ! (mysql_exec "${replica}" "SELECT 1") >/dev/null; do
    echo "Waiting for replica database ${replica} connection..."
    sleep 5
  done

  replication_status=$(mysql_exec "${replica}" "SELECT * FROM performance_schema.replication_group_members WHERE member_host <> NULL")
  if [[ -n "${replication_status}" ]]; then
    echo "Replication already started:"
    echo "${replication_status}"
    exit 0
  fi
done


echo "creating replica user"
mysql_exec db_1 "SET SQL_LOG_BIN=0;"
mysql_exec db_1 "CREATE USER rpl_user@'%' IDENTIFIED BY 'password';"
mysql_exec db_1 "GRANT REPLICATION SLAVE ON *.* TO rpl_user@'%';"
mysql_exec db_1 "GRANT CONNECTION_ADMIN ON *.* TO rpl_user@'%';"
mysql_exec db_1 "GRANT BACKUP_ADMIN ON *.* TO rpl_user@'%';"
mysql_exec db_1 "GRANT GROUP_REPLICATION_STREAM ON *.* TO rpl_user@'%';"
mysql_exec db_1 "CREATE USER 'haproxy_check'@'%';"
mysql_exec db_1 "FLUSH PRIVILEGES;"
mysql_exec db_1 "SET SQL_LOG_BIN=1;"

echo "configure replication source"
mysql_exec db_1 db_2 db_3 "CHANGE REPLICATION SOURCE TO SOURCE_USER='rpl_user', SOURCE_PASSWORD='password' FOR CHANNEL 'group_replication_recovery';"

echo "bootstrapping"
mysql_exec db_1 "SET GLOBAL group_replication_bootstrap_group=ON;"
mysql_exec db_1 "START GROUP_REPLICATION USER='rpl_user', PASSWORD='password';"
mysql_exec db_1 "SET GLOBAL group_replication_bootstrap_group=OFF;"

mysql_exec db_2 db_3 "START GROUP_REPLICATION USER='rpl_user', PASSWORD='password';"

mysql_exec db_1 "SELECT * FROM performance_schema.replication_group_members;"

echo "creating schema"
for f in /docker-entrypoint-initdb.d/* ; do
  schema=$(cat "${f}")
  mysql -h db_1 --user=root --password="${MYSQL_ROOT_PASSWORD:-password}" -D db -e "${schema}";
done