FROM golang:latest

WORKDIR /home/accounting/
COPY ./server .
RUN /bin/bash ./installRequirements
COPY ./accounting.conf /etc/accounting.conf
RUN mkdir /var/log/accounting/

CMD go run . -c /etc/accounting.conf
