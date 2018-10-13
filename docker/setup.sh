#!/bin/bash

###########THIS SETUP DOES WORK PROPERLY, NEEED TO BE FIXED

#MONGODB1=`ping -c 1 mongo1 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`
#MONGODB2=`ping -c 1 mongo2 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`
#MONGODB3=`ping -c 1 mongo3 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`

MONGODB1="localhost:27027"
MONGODB2="localhost:27028"
MONGODB3="localhost:27029"

echo "Started.."

echo SETUP.sh time now: `date +"%T" `
mongo --host ${MONGODB1} <<EOF
use admin
rsconf = {
  _id: "rs",
  members: [
    {
     _id: 0,
     host: "mongodb-main:27017"
    },
    {
     _id: 1,
     host: "mongodb-sub1:27017"
    },
    {
     _id: 2,
     host: "mongodb-sub2:27017"
    }
   ]
};
rs.initiate(rsconf);
rs.reconfig(rsconf);
rs.slaveOk();
db.getMongo().setReadPref('nearest');
db.getMongo().setSlaveOk(); 
EOF
