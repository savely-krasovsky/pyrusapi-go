# pyrusapi-go

[![GoDoc Widget](https://godoc.org/github.com/L11R/pyrusapi-go?status.svg)](https://godoc.org/github.com/L11R/pyrusapi-go)
[![Go Report](https://goreportcard.com/badge/github.com/L11R/pyrusapi-go)](https://goreportcard.com/report/github.com/L11R/pyrusapi-go)
[![codecov](https://codecov.io/gh/L11R/pyrusapi-go/branch/master/graph/badge.svg)](https://codecov.io/gh/L11R/pyrusapi-go)
![test](https://github.com/L11R/pyrusapi-go/actions/workflows/test.yml/badge.svg)
![lint](https://github.com/L11R/pyrusapi-go/actions/workflows/lint.yml/badge.svg)


Library to work with Pyrus API v4 written in Golang.

## Install
`go get github.com/L11R/pyrusapi-go`

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/L11R/pyrusapi-go"
)

func main() {
	c, err := pyrus.NewClient("bot_login", "bot_security_key")
	if err != nil {
		log.Fatalln(err)
	}

	p, err := c.Profile()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(p.FirstName)
}
```

## Current status

Forms:
- [x] `GET /forms`
- [x] `GET /forms/{form-id}`
- [x] `GET /forms/{form-id}/register`

Tasks:
- [x] `GET /tasks/{task-id}`
- [x] `POST /tasks`
- [x] `POST /tasks/{task-id}/comments`

Files:
- [x] `POST /files/upload`
- [x] `GET /files/download/{file-id}`

Catalogs:
- [x] `GET /catalogs`
- [x] `GET /catalogs/{catalog-id}`
- [x] `PUT /catalogs`
- [x] `POST /catalogs/{catalog-id}`

Contacts:
- [x] `GET /contacts`

Members:
- [x] `GET /members`
- [x] `POST /members`
- [x] `PUT /members/{member-id}`
- [x] `DELETE /members/{member-id}`

Lists:
- [x] `GET /lists`
- [x] `GET /lists/{list-id}/tasks`
- [x] `GET /inbox`

Telephony:
- [x] `GET /calls`
- [x] `PUT /calls/{call-guid}`
- [x] `POST /calls/{call-guid}/event`

Webhooks:
- [x] Use `WebhookHandler() (http.HandlerFunc, <-chan Event)`

# Tests
To test project you have to download sample responses from your organization. I cannot upload my own, because it could
contain sensitive information, and I'm too lazy to fake it :)

Set those environment variables and run tests:
```dotenv
PYRUS_LOGIN=login
PYRUS_SECURITY_KEY=securityKey
PYRUS_FORM_ID=123456
PYRUS_TASK_ID=123456
PYRUS_CATALOG_ID=123456
PYRUS_MEMBER_ID=123456
PYRUS_ROLE_ID=123456
PYRUS_LIST_ID=123456
PYRUS_FILE_ID=123456
PYRUS_CALL_GUID=5d8dc3d6-27e7-4cd4-a057-2b4f4d74e0a5
```