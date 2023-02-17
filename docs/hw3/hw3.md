# Асинхронная репликация

1. Переводим запросы /user/get/{id} и /user/search на чтение с реплики
2. Запускаем 
   ```
   docker-compose -f docker-compose_async_repl.yml up
   ```
   
3. Настраиваем репликацию скриптом
   ```
   db/async/replication_init.sh
   ```

4. Создаем схему бд и наполняем таблицу

   ```
   docker exec -i highload_architect-db-1 mysql -t -uroot -ppassword db < db/01_schema.sql
   docker exec -i highload_architect-db-1 mysql -t -uroot -ppassword db < db/02_load_people.sql
   ```

5. Запускаем jmeter c нагрузкой на /user/get/{id} и /user/search ``jmeter -t test/Test\ Plan.jmx``
6. Проверяем потребление ресурсов ``docker stats``
```
CONTAINER ID   NAME                              CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O         PIDS
8f9dd22642b1   highload_architect-app-1          1.96%     12.41MiB / 7.773GiB   0.16%     46.1MB / 67.2MB   5.46MB / 0B       7
f06b9684ecb6   highload_architect-db_replica-1   133.48%   505.4MiB / 7.773GiB   6.35%     1.33MB / 45.2MB   33.7MB / 15.2MB   39
dbd3cb1b21e8   highload_architect-db-1           1.16%     359.2MiB / 7.773GiB   4.51%     1.17kB / 0B       12.4MB / 15.2MB   38
```

# Полусинхронная репликация
1. Запускаем контейнеры 
   ```
   docker-compose -f docker-compose_semisync_repl.yml up
   ```
   В файлах конфигурации настроены GTID и Row based replication `db/semisync/mysql.cnf` `db/semisync/mysql_replica_1.cnf` `db/semisync/mysql_replica_2.cnf`

2. Настраиваем semisync репликацию 
   ```
   db/semisync/replication_init.sh
   ```
   
3. Создаем схему бд 
    
   ```
   docker exec -i highload_architect-db-1 mysql -t -uroot -ppassword db < db/01_schema.sql
   ```
   
4. Запускаем запись в таблицу 

   ```
   db/semisync/write_into_users.sh highload_architect-db-1
   ``` 
   В выводе будет количество завершенных транзакций

5. Останавливаем мастер БД 
   ```
   docker kill highload_architect-db-1
   ```
   
6. Проверяем количество записей в репликах
   ```
   docker exec highload_architect-db_replica_1-1 mysql -uroot -ppassword db -e "select count(*) from users"
   docker exec highload_architect-db_replica_2-1 mysql -uroot -ppassword db -e "select count(*) from users"
   ```
   
7. Продвигаем одну из оставшихся реплик до мастера 
   ```
   db/semisync/change_source.sh
   ```
   
9. Вставляем запись в БД и проверяем что запись реплицировалась
    ```
   docker exec highload_architect-db_replica_1-1 mysql -uroot -ppassword db -e "INSERT INTO db.users (first_name, second_name, password_hash) VALUES ('Ivan', 'Petrov', '12345')"
   docker exec highload_architect-db_replica_1-1 mysql -uroot -ppassword db -e "select count(*) from users"
   docker exec highload_architect-db_replica_2-1 mysql -uroot -ppassword db -e "select count(*) from users"
   ```
