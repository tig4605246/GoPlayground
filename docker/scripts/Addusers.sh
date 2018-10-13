#!/bin/bash

MONGODB1="localhost:27027"

mongo --host ${MONGODB1} <<EOF
use admin
kelier = {
user: "kelier",
pwd: "123",
roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
};

kelly = {
user: "kelly",
pwd: "123",
roles: [ { role: "root", db: "admin" } ]
};
db.createUser(kelier);
db.createUser(kelly);
db.getUsers();
EOF