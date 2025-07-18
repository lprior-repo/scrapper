# Environment-specific configurations for different test environments

# Development environment
[development]
variable "base_url" {
    default = "http://localhost:8080"
}

variable "github_token" {
    default = ""
}

variable "neo4j_uri" {
    default = "bolt://localhost:7687"
}

variable "debug_mode" {
    default = true
}

# Testing environment
[testing]
variable "base_url" {
    default = "http://localhost:8081"
}

variable "github_token" {
    default = ""
}

variable "neo4j_uri" {
    default = "bolt://localhost:7687"
}

variable "debug_mode" {
    default = true
}

# Staging environment
[staging]
variable "base_url" {
    default = "https://staging-api.example.com"
}

variable "github_token" {
    default = ""
}

variable "neo4j_uri" {
    default = "bolt://staging-neo4j.example.com:7687"
}

variable "debug_mode" {
    default = false
}

# Production environment (read-only tests)
[production]
variable "base_url" {
    default = "https://api.example.com"
}

variable "github_token" {
    default = ""
}

variable "neo4j_uri" {
    default = "bolt://prod-neo4j.example.com:7687"
}

variable "debug_mode" {
    default = false
}

# CI/CD environment
[ci]
variable "base_url" {
    default = "http://localhost:8080"
}

variable "github_token" {
    default = ""
}

variable "neo4j_uri" {
    default = "bolt://localhost:7687"
}

variable "debug_mode" {
    default = false
}

variable "test_timeout" {
    default = 60
}

# Docker environment
[docker]
variable "base_url" {
    default = "http://api:8080"
}

variable "github_token" {
    default = ""
}

variable "neo4j_uri" {
    default = "bolt://neo4j:7687"
}

variable "debug_mode" {
    default = false
}