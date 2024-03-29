# Enterprise Note Application

## Team members
* Floyd Watson
* Travis Falkenberg

## Deploy

to deploy as dev through Docker CLI

### Client standalone:

Testing version

        docker build -f Dockerfile.dev -t my-golang-app .
        docker run -p 8080:8080 -it --rm --name my-running-app my-golang-app  

Application will not be able to talk to database in this state making database.go tests fail

---
Running Version

        docker build -t my-golang-app .
        docker run -p 8080:8080 -it --rm --name my-running-app my-golang-app  

Home page can be reached at < docker ip address >:8080

This state cannot talk to database so will cause errors when sdatabase is called

---


### Full local deployment using docker-compose

Start database first
Make sure to change connection string in database.go to provided connString for docker connection

        docker-compose up -d db
        docker-compose logs -f db      

And look out for a log line like:

        db_1   | LOG:  database system is ready to accept connections

Then start web app with 

        docker-compose up web   



## Current Issues

Currently cannot fully test docker-compose as using windows home and that operating system cannot currently make use of volumes

---

### Deploy as local exe

1. A PostgreSQL database must be created called ***noteBookApp***

![Imgur](https://i.imgur.com/nKJrXbr.png)

2. Navigate to init.sql and copy all of the contents

![Imgur](https://i.imgur.com/Efpk27p.png)

3. Create a query from Database created and copy SQL into query. Run query

![Imgur](https://i.imgur.com/evZXSDG.png)

4. Navigate to client folder and run go build

![Imgur](https://i.imgur.com/RZ8nA9P.png)

5. Navigate to containing folder and run client.exe

![Imgur](https://i.imgur.com/iFGBZhk.png)

6. Open your browser and navigate to localhost:8080

![Imgur](https://i.imgur.com/AVBuTAQ.png)

7. Login using preset user. Or create your own. Preset user details, username: Trav3, password: 1234

![Imgur](https://i.imgur.com/otfE0Qr.png)
