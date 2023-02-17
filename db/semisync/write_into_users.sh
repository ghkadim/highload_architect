#!/bin/bash
set -x

db=$1

docker exec -i "${db}" /bin/bash <<'EOF'
i=0
var=foo
while true; do
  mysql --user=root --password=password -e \
    "INSERT INTO db.users (first_name, second_name, password_hash) VALUES ('Ivan-${i}', 'Petrov-${i}', '12345')" > /dev/null 2>&1 || exit 1
  i=$((i + 1))
  echo "Transaction: ${i}"
done
EOF
