FROM proxysql/proxysql:2.5.2

RUN apt update -y --fix-missing && apt install -y mariadb-client

ADD proxysql.cnf /etc/