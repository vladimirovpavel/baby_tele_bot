from docker.io/library/postgres:13.6
VOLUME /baby_postgres
EXPOSE 5432/tcp
COPY init.sql /docker-entrypoint-initdb.d/
#--mount type=volume,src=baby_postgres,dst=/var/lib/postgresql/data
