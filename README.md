# Rev Thru Rice Backend
Name: Lim Jinle, Alexius

## Overview
Rev Thru Rice is a forum designed for travellers in Southeast Asia who enjoy exploring this beautiful region on two bikes.

API functions: Authentication; Creation, update and deletion of user, post, review, comment and rating

Frontend: https://github.com/TastyVeggy/rev-thru-rice-frontend

## Setup

### Docker

1. Make a copy of `.exampleenv` and name the copy `.env` and change the necessary settings. You should not need to change the `DATABASE_URL`

2. Set up and run a containerised local deployment consisting of both the backend server and the database. You should run the following command in the same folder as your .env file

   ```
   docker compose up -d
   ```

### Manually

#### PostgreSQL

You can find a host provider and host the database online or host it locally by downloading PostgreSQL on your local machine or use a docker image.

Create a database to be used by the server

#### Configuration

1. Make a copy of `.exampleenv` and name the copy `.env`

2. Set `DB_URL` to `postgres://<user>:<password>@<ip>:<port>/<database>`

3. Set `JWT_SECRET_KEY` to your own secret key (signing method is HS256), you can generate it with OpenSSL

4. Running for the first time, you should set `CREATE_TABLE` to `TRUE` to do the necessary migrations but can be set to false in subsequent runs though not necessary. `SEED_DATA_DIR` is the directory for seeding the initial data and should be left to the default on the first run and can be set to nothing in subsequent runs though not necessary.

#### Running the backend

1. Install necessary dependencies
   ```
   go mod tidy
   ```
2. Run the application
   ```
   go run server.go
   ```
