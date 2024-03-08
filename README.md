# Project Name

## Overview

This project is a web application backend implemented in Go, using the `mux` router for efficient HTTP request routing. It features a RESTful API design, supporting various endpoints for operations such as user authentication, account management, and content manipulation. Middleware functions are used for authorization, CORS handling, and admin verification, ensuring secure and flexible interaction with the frontend.

## Features

- User authentication and authorization with JWT
- CRUD operations on user accounts and posts
- Admin-specific operations for enhanced control
- CORS support for cross-origin requests
- Extensive use of middleware for request handling

## API Endpoints
- GET /: Home page
- GET /band/{id}: Retrieve band information by ID
- POST /delete/{id}: Delete a user by ID (authorized JWT users only)
- PUT /update/{id}: Update a user by ID (authorized JWT users only)
- GET /admin: Admin panel (admins only)
- POST /login: Login endpoint
- POST /signup: Signup endpoint
- GET /posts: Retrieve all posts (GET and OPTIONS methods)
- POST /posts: Create a new post (POST and OPTIONS methods)
- GET /posts/{id}: Retrieve a post by ID (GET and OPTIONS methods)
- PUT /posts/{id}: Update a post by ID (PUT and OPTIONS methods)
- DELETE /posts/{id}: Delete a post by ID (DELETE and OPTIONS methods)