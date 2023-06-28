# Запуск приложения с 2 шардами
```bash
docker-compose -f docker-compose_sharding.yml up
```

# Наполняем БД диалогами
```bash
TEST_BEFORE_RESHARDING=1 pytest -v test/api/test_dialog.py::test_multiple_dialogs_create
```

# Проверяем шарды
```bash
mysql -u user -ppassword -P3361 -h 127.0.0.1
mysql -u user -ppassword -P3362 -h 127.0.0.1
mysql -u user -ppassword -P3363 -h 127.0.0.1
```

```sql
select * from db.dialogs;
```

Подключемся к proxysql, проверяем роутинг
```bash
mysql -u user -ppassword -P6032 -h 127.0.0.1
```

```sql
select active,hits, mysql_query_rules.rule_id, match_digest, match_pattern, replace_pattern, cache_ttl, apply,flagIn,flagOUT FROM mysql_query_rules NATURAL JOIN stats.stats_mysql_query_rules ORDER BY mysql_query_rules.rule_id;
```

# Подготовка 3 шарды, решардинг
```bash
go run cmd/reshard/main.go -config reshard.conf.yaml -debug
```

проверяем что в 3 шарду скопированы данные
```bash
mysql -u user -ppassword -P3363 -h 127.0.0.1
```

```sql
select * from db.dialogs;
```

# Настройка proxysql

```bash
mysql -u user -ppassword -P6032 -h 127.0.0.1
```

Добавляем новый инстанс
```sql
INSERT INTO mysql_servers (hostgroup_id,hostname,port,max_connections) VALUES (2,'db_3',3306,200);
LOAD MYSQL SERVERS TO RUNTIME; SAVE MYSQL SERVERS TO DISK;
```

Обновляем правила шардирования
```sql
DELETE FROM mysql_query_rules WHERE rule_id IN (1, 2);
INSERT INTO mysql_query_rules (rule_id,active,match_pattern,destination_hostgroup,apply) VALUES (1,1,"(00\d|0[1-9]\d|[12]\d{2}|3[0-2]\d|33[0-3])\d*",0,1);
INSERT INTO mysql_query_rules (rule_id,active,match_pattern,destination_hostgroup,apply) VALUES (2,1,"(33[4-9]|3[4-9]\d|[45]\d{2}|6[0-5]\d|66[0-6])\d*",1,1);
INSERT INTO mysql_query_rules (rule_id,active,match_pattern,destination_hostgroup,apply) VALUES (3,1,"(66[7-9]|6[7-9]\d|[7-9]\d{2})\d*",2,1);
LOAD MYSQL QUERY RULES TO RUNTIME;SAVE MYSQL QUERY RULES TO DISK;
```

# Проверяем чтение диалогов
```bash
TEST_AFTER_RESHARDING=1 pytest -v test/api/test_dialog.py::test_multiple_dialogs_read
```
Подключемся к proxysql, проверяем роутинг
```bash
mysql -u user -ppassword -P6032 -h 127.0.0.1
```

```sql
select active,hits, mysql_query_rules.rule_id, match_digest, match_pattern, replace_pattern, cache_ttl, apply,flagIn,flagOUT FROM mysql_query_rules NATURAL JOIN stats.stats_mysql_query_rules ORDER BY mysql_query_rules.rule_id;
```

# Очищаем шарды от лишних данных, проверяем
```bash
go run cmd/reshard/main.go -config reshard.conf.yaml -debug -mode cleanup
```

```bash
mysql -u user -ppassword -P3361 -h 127.0.0.1
mysql -u user -ppassword -P3362 -h 127.0.0.1
mysql -u user -ppassword -P3363 -h 127.0.0.1
```

```sql
select * from db.dialogs;
```

# Останавливаем приложени, очищаем
```bash
docker-compose -f docker-compose_sharding.yml down
make compose_clean
```
