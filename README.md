
### Folder structure
```md
│
├── cmd
│   └── server
│       └── main.go
│
├── config
│   └── config.go
│
├── internal
│
│   ├── api
│   │   ├── routes
│   │   │   ├── product.routes.go
│   │   │   ├── cart.routes.go
│   │   │   └── order.routes.go
│   │   │
│   │   └── handlers
│   │       ├── product.handler.go
│   │       ├── cart.handler.go
│   │       └── order.handler.go
│   │
│   ├── services
│   │   ├── product.service.go
│   │   ├── cart.service.go
│   │   └── order.service.go
│   │
│   ├── repository
│   │   ├── product.repo.go
│   │   ├── cart.repo.go
│   │   └── order.repo.go
│   │
│   ├── models
│   │   ├── product.model.go
│   │   ├── cart.model.go
│   │   └── order.model.go
│   │
│   ├── middleware
│   │   ├── logger.go
│   │   ├── errorHandler.go
│   │   └── requestID.go
│   │
│   ├── dto
│   │   ├── product.dto.go
│   │   ├── cart.dto.go
│   │   └── order.dto.go
│   │
│   └── utils
│       ├── response.go
│       ├── validator.go
│       └── discount.go
│
├── pkg
│   └── database
│       └── mongo.go
│
├── scripts
│   └── seed.go
│
├── docs
│   └── openapi.yaml
│
├── .env
├── go.mod
└── README.md
```