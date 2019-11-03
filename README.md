"# EnterpriseNoteApp" 

---

## Deploy

to deploy as dev through Docker CLI

### Client standalone:

        docker build -f Dockerfile.dev -t my-golang-app .
        docker run -p 8080:8080 -it --rm --name my-running-app my-golang-app  

Application will not be able to talk to database in this state. this is just used for testing purposes

Home page can be reached at < docker ip address >:8080


### Full local deployment using docker-compose

Start database first

        docker-compose up -d db
        docker-compose logs -f db      

And look out for a log line like:

        db_1   | LOG:  database system is ready to accept connections

Then start web app with 

        docker-compose up web   



## Current Issues

DB conn string is not correct, trying to find error
