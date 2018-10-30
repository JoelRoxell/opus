FROM centos

RUN mkdir -p /opt/build/

WORKDIR /opt/build

COPY test.txt .
ADD data storage

CMD ["tail", "-f", "/dev/null"]
