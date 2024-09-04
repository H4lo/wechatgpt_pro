#!/bin/bash

service mariadb start
echo "CREATE USER 'test'@'%' IDENTIFIED BY 'test'; GRANT ALL PRIVILEGES ON *.* TO 'test'@'%'; FLUSH PRIVILEGES;" > /tmp/mysql-setup.sql
mysql < /tmp/mysql-setup.sql

mysql < /root/build/db/init.sql

/root/build/server
