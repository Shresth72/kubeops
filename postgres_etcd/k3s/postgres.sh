#!/bin/bash

# Create a new user using the root user and grant given permission
# (Only required for Kine in Postgres 15)
sudo -u postgres psql

CREATE DATABASE kubernetes;

CREATE USER shres WITH ENCRYPTED PASSWORD 'secret';

GRANT ALL PRIVILEGES ON DATABASE kubernetes TO shres;

# Connect to kubernetes db from root
\c kubernetes postgres

GRANT ALL ON SCHEMA public TO shres;


