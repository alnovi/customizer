FROM mongo:latest

# COPY mongo.conf /usr/local/etc/mongo/mongo.conf
# COPY mongo-init.sh /docker-entrypoint-initdb.d/mongo-init.sh

VOLUME /data/db /data/configdb

CMD ["mongod"]

EXPOSE 27017