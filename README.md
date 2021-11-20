# To-do
- Get player stat data from game
    - There should be a versioning system. Current values could be saved as current?
- Create models for balance values. Maybe a big json with all game values inside? 
- A way to get client Id for the game. Because each client's current balance version will be different.


# Done
## 2021-19-12
- Add http client code in game 
- Download & Install postman
- Fix backend asking firewall permission every time it is launched
- Figure out how to parse request body
- Reorganize folders to allow for api versioning.
    - Forgot to call structs with `[package].[struct]` wasted like 10 minutes again. One day I will learn.
- Create versions table

## 2021-19-11
- Installed sqlx package for backend
- Installed pgx package for backend
- Installed PostgreSQL
- Created database Hyperion (Codename for our backend)
- Created table users (Had to learn how to create table from file using psql)
- Add auto increment id
- Add created_date default value
- I guess we don't need sqlx when we have pgx?
- Learn how to connect to database using pgx
- Select values from database using pgx
- Remove unused packages

