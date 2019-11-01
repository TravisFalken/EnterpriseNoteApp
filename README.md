"# EnterpriseNoteApp" 

---

## Deploy

to deploy as dev

Client standalone:

        docker build -f Dockerfile.dev -t my-golang-app .
        docker run -p 8080:8080 -it --rm --name my-running-app my-golang-app  
