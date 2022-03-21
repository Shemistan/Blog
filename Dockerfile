FROM postgres:13.3

LABEL description = "postgresql instance"/

COPY migrations/init.sql /docker-entrypoint-initdb.d/
