# USDC Transfer Listener Service

## Description

This service listens Transfer events from USDC smart-contract on Ethereum mainnet and stores information about transactions to PostgreSQL database

## Install

  ```
  git clone github.com/bohdan-vykhovanets/usdc-transfer-listener-svc
  cd usdc-transfer-listener-svc
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main migrate up
  ./main run service
  ```

## Running from docker 
  
1. Make sure that docker installed.
2. Download docker-compose.yaml
3. Add your actual Infura API key to .env file as INFURA_API so docker-compose can read it.
4. Execute `docker compose up`

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `migrate up` command to create database schema
* Launch the service with `run service` command


### Database
For services, we do use ***PostgresSQL*** database. 
You can [install it locally](https://www.postgresql.org/download/) or use [docker image](https://hub.docker.com/_/postgres/).
