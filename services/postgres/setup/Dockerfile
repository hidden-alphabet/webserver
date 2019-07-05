FROM postgres

RUN apt update -y && \
    apt install -y git

COPY init.sh ./
COPY *.sql ./

ENTRYPOINT ["/bin/bash", "./init.sh"]