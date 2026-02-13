# 🧠 Redis Clone — In-Memory Key-Value Store in Go

A simple *Redis Clone built in Go* — great for learning how Redis works under the hood and experimenting with custom data stores.

This project implements core parts of the Redis protocol and command set using Go’s network and concurrency primitives. It’s designed for educational purposes, performance experimentation, and as a stepping stone to deeper understanding of databases and distributed systems.

👉 Check out an article create by myself to learn how to do your own: https://alejandro.buzz/projects/redis

🚀 Features

- ✔️ RESP (Redis Serialization Protocol) implementation.
- ✔️ In-memory data storage with simple data types
- ✔️ Supports concurrent clients via goroutines.
- ✔️ Custom protocol parser written in Go.

(Add more features as you implement them!)

# 📌 Why This Project Exists

Redis is a high-performance in-memory database widely used for caching, messaging, session stores, and analytics. This project helps you:

- Understand how Redis parses commands (RESP).
- Explore building a database server from scratch in Go.
- Learn about TCP networking, concurrency, and protocol design.

(It’s a learning project, not a production-ready database.)

# 💾 Installation
1. **Clone the Repository**
```bash
$ git clone https://github.com/dyxgou/redis
```

2. Download the needed dependencies (just godotenv).

```bash
$ go mod download
```

3. **Building and Running the Redis Server**

Compile the server using the provided Makefile:

```bash
$ make
```

4. **Build the Redis Client**
To compile the custom Redis client:

```bash
$ make client
```
