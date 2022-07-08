#! /bin/sh
mongo --port 27017 -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD --authenticationDatabase admin <<EOF
use $MONGO_CUSTOMIZER_DATABASE
db.createUser(
    {
        user: "$MONGO_CUSTOMIZER_USER",
        pwd: "$MONGO_CUSTOMIZER_PASSWORD",
        roles: [
            {
                role: "dbOwner",
                db: "$MONGO_CUSTOMIZER_DATABASE"
            }
        ],
        mechanisms: [ "SCRAM-SHA-1" ]
    }
);
EOF