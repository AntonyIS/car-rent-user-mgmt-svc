# Notlify User Service
This is part of the Notlify application that manages users and user data.
* Responsible for user registration, authentication, and authorization.
* Stores user profile information, handles user preferences, and manages user roles.
* Provides endpoints for user-related operations like sign up, login, profile updates, etc

## Features
1. User authentication and authorization
    * Implement role-based access control to manage user permissions and privileges.
    * Define user roles (e.g., regular user, author, moderator) to control access to specific features.
2. User profile management
    * Enable users to update their profile information, such as name, bio, profile picture, and contact details.
    * Allow users to customize their settings and preferences.
3. Password reset and recovery
    * Implement a secure process for users to reset their passwords if forgotten.
    * Send password reset links via email with one-time tokens for verification.

## System design components
1. Cloud services (AWS)
    * Databases
        1. PostgreSQL
        2. S3
    * Deployment
        1. Docker (ECR)
        2. ECS
        3. EKS
    * Hosting 
        1. Route53

2. Technology stack
    * Go programming laguages
    * Go Gin framework
    * Docker 
    * Kubernetes
    * ReactJS + Typescript
    * AWS 

## Application structure
This service uses Hexagonal architecture
```
.
├── config
│   └── config.go
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── app
│   │   │   ├── controllers.go
│   │   │   ├── controllers_test.go
│   │   │   └── handler.go
│   │   └── repository
│   │       ├── mongodb
│   │       │   └── mongodb.go
│   │       ├── postgres
│   │       │   └── postgres.go
│   │       └── s3
│   │           └── s3.go
│   └── core
│       ├── domain
│       │   └── domain.go
│       ├── ports
│       │   └── ports.go
│       └── services
│           └── services.go
├── LICENSE
├── main.go
├── Makefile
└── README.md
```


