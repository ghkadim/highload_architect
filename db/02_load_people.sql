SET GLOBAL max_connections = 1000;

LOAD DATA INFILE '/var/lib/mysql-files/people.csv'
  INTO TABLE users
  FIELDS TERMINATED BY ','
  LINES TERMINATED BY '\n'
  (@name, age, city)
  SET
    second_name = SUBSTRING_INDEX(@name, ' ', 1),
    first_name = SUBSTRING_INDEX(@name, ' ', -1),
    password_hash = CAST('password' AS BINARY);