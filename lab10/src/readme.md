поменять в /etc/postgresql/14/main/postgresql.conf строку listen_addresses (localhost по умолч.) на '*'

поменять в /etc/postgresql/14/main/pg_hba.conf строку:  
host    all             postgres             0.0.0.0/0               md5  
на  
host    all             all             0.0.0.0/0               md5

Итоговые настройки pg_hba.conf:
# Database administrative login by Unix domain socket
local   all             postgres                                md5

# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all                                     peer
# IPv4 local connections:
host    all             all             127.0.0.1/32            trust
# IPv6 local connections:
host    all             all             ::1/128                 trust
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     peer
host    replication     all             127.0.0.1/32            scram-sha-256
host    replication     all             ::1/128                 scram-sha-256
#Allow connections from Docker container
host    bmstu_db        postgres        172.18.0.2/32           md5