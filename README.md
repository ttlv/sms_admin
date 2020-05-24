# 1.What is it?

This is a CRM Server for sms server

## https://github.com/ttlv/sms

# 2. How to run it

## 1. Initialize Env

### 1 .DB config env

````
export SMSADMIN_DBNAME=sms_dev
export SMSADMIN_HOST=127.0.0.1
export SMSADMIN_DBPort=7100
export SMSADMIN_User=sms
export SMSADMIN_Password=123
````



### 2 . SMSADMIN Config

````
export SMSADMIN_SERVERPORT=7000
export SMSADMIN_HTTPS=false
export SMSADMIN_HTTPAUTHNAME=sms
export SMSADMIN_HTTPAUTHPASSWORD=sms
````

## 2. Use Docker To Run Dependent Services

````
docker-compose -f docker-compose.yml up -d
````

## 3. Run Admin Server

````
go run main.go -compile-templates=true
````

## 3. Have a look

input http://localhost:7000/admin


