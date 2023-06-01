FROM mysql:8.0

#RUN mkdir /usr/local/shell
#COPY ./docker/db/shell /usr/local/shell

#RUN echo "now building..."

#RUN mkdir docker-entrypoint-initdb.d
#COPY ./docker/db/sql ./docker-entrypoint-initdb.d/
#RUN chmod 0775 docker-entrypoint-initdb.d/init.sh
#RUN /bin/bash -c "./docker-entrypoint-initdb.d/init.sh"