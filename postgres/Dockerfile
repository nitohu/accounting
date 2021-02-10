FROM postgres:12.2

COPY ./init/1-init.sql /docker-entrypoint-initdb.d/
COPY ./init/2-tables.sql /docker-entrypoint-initdb.d/
COPY ./init/3-data.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_PASSWORD "supersecret"
EXPOSE 5432
