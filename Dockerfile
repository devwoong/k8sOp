FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install vim -y

RUN mkdir /home/server

WORKDIR /home/server

COPY k8sOp .

ADD apply /home/server/apply

ENTRYPOINT [ "/home/server/k8sOp" ]

# EXPOSE 80