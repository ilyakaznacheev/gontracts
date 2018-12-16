# Simple contract microservice

The microservice able to handle companies, contracts, and purchases

[![Go Report Card](https://goreportcard.com/badge/github.com/ilyakaznacheev/gontracts)](https://goreportcard.com/report/github.com/Masterminds/glide) [![GoDoc](https://godoc.org/github.com/ilyakaznacheev/gontracts?status.svg)](https://godoc.org/github.com/ilyakaznacheev/gontracts) [![Build Status](https://travis-ci.org/ilyakaznacheev/gontracts.svg?branch=master)](https://travis-ci.org/ilyakaznacheev/gontracts) [![Coverage Status](https://coveralls.io/repos/github/ilyakaznacheev/gontracts/badge.svg?branch=master)](https://coveralls.io/github/ilyakaznacheev/gontracts?branch=master)

## Usage

To run service you have to run MySQL server via Docker

```shell
docker-compose up
```

and then start the server

```shell
go run example/run.go
```

The server will be started on `localhost:8000`

## Requirements

If you download the package manually following requirements must be met:
- github.com/go-sql-driver/mysql
- github.com/gorilla/mux

## API

Service has following web API:
- `/company/<id:int>` GET: get company data by id 
- `/company` POST: create new company
- `/company` PUT: create or update company
- `/company/<id:int>` DELETE: delete company by id
- `/companies` GET: get list of all companies
- `/contract/<id:int>` GET: get contract data by id 
- `/contract` POST: create new contract
- `/contract` PUT: create or update contract
- `/contract/<id:int>` DELETE: delete contract by id
- `/contracts` GET: get list of all contracts
- `/contract/<id:int>/purchase` GET: get purchase history of contract
- `/purchase` POST: create new purchase document

## Examples

### Get company data

**Request**

GET:`localhost:8000/company/1`

**Response**

```json
{
    "ID": 1,
    "name": "Megacom",
    "regcode": "MGC111"
}
```

### Create a company

**Request**

POST:`localhost:8000/company`
```json
{
	"name":"Supercom",
	"regcode":"SRC222"
}
```

**Response**

```json
{
	"ID": 2
}
```

### Update company info

**Request**

PUT:`localhost:8000/company`
```json
{
	"name":"Supercom",
	"regcode":"SRC333"
}
```

**Response**

```json
{
	"name":"Supercom",
	"regcode":"SRC333"
}
```
### Delete company

**Request**

DELETE:`localhost:8000/company/3`

**Response**

`OK`

### Add new contract

**Request**

POST:`localhost:8000/contract`
```json
{
	"sellerID":1,
	"clientID":2,
	"validFrom":"2000-01-01T00:00:00Z",
	"validTo":"2001-01-01T00:00:00Z",
	"amount":150
}
```

**Response**

```json
{
	"ID": 5
}
```

### Add new purchase document

**Request**

POST:`localhost:8000/purchase`
```json
{
	"contractID":1,
	"datetime": "2000-10-01T00:00:00Z",
	"amount": 3
}
```

**Response**

```json
{
	"ID": 4
}
```

### Get purchase history of contract

**Request**

GET:`localhost:8000/contract/1/purchase`

**Response**

```json
[
	{
		"ID": 1,
		"contractID": 1,
		"datetime": "2000-10-01T00:00:00Z",
		"amount": 3
	},
	{
		"ID": 2,
		"contractID": 1,
		"datetime": "2000-10-01T00:00:00Z",
		"amount": 3
	},
	{
		"ID": 3,
		"contractID": 1,
		"datetime": "2000-10-01T00:00:00Z",
		"amount": 3
	},
	{
		"ID": 4,
		"contractID": 1,
		"datetime": "2000-10-01T00:00:00Z",
		"amount": 3
	}
]
```