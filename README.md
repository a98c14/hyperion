# To-do
- Get player stat data from game
    - There should be a versioning system. Current values could be saved as current?
- Create models for balance values. Maybe a big json with all game values inside? 
- A way to get client Id for the game. Because each client's current balance version will be different.
- [DevTools] List all possible spawnable entities when /help is used

# Done
## 2021-11-24
- Create database design docs
- Design database for dynamic prefab system
- Create tables 

## 2021-11-23
- Find a way to dynamically load entities to UnityEngine
## 2021-11-22
- Create a vscode task for build and run. Add a shortcut for restarting currently running task
- Brainstorm about ideas on how to create an editor that allows us to ditch gameobject conversion entirely
## 2021-11-20
- Add http client code in game 
- Download & Install postman
- Fix backend asking firewall permission every time it is launched
- Figure out how to parse request body
- Reorganize folders to allow for api versioning.
    - Forgot to call structs with `[package].[struct]` wasted like 10 minutes again. One day I will learn.
- Create versions table
- Rename go module to hyperion
- Add versions handler
    - Get all function
- Fix import cycling
- Fix compilation errors
- Add shortcuts for `toggle breakpoint` and `next error in files`
- Learn how to debug go apps in vs code
- Add more shortcuts for debugging
- Fix runtime errors for version request
## 2021-11-19
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

