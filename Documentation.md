# Enterprise Note Application Documentation

## Deployment Requirements

### Client side candidates

* We have been using the client side for storing a session ID cookie. This is kept on the client side so when ever the client sends a request to the server the cookie ID will be sent with the request so we know who is sending the request and if client is in fact “logged on” The pros of this is that it is easy for the server to Identify which request belongs to who and the cons are that if anyone else gets the value of the cookie and “impersonates” the client they can access the clients information.


### Server side candidates

* We are saving all of the notes the user creates on the server side because it is more secure and easier to query than having to get it from the client every time.
* We have been utilizing the server side to do our searches through SQL. We use this technique to increase security for other users notes. If we searched through these notes on client side then it would mean having to send all notes through to client for them to be searched through. This allows for people without access to those notes to be able to intercept them more easily.
* When the client creates an account we are saving the clients details along with their username and password on the server side. This makes their details more secure and when querying the details with sql and managing details a lot easier. The downside is that the company is responsible for the sensitive data and if any data is leaked the company takes responsibility. 

---

### Database design choices

* We initially wanted to make a CRUD like database interface. In future we would be removing alot of the repeated code in our current database.go file and creating a much more CRUD system.

---

### Deploying on different operating systems

* We are utilizing Docker to host our application. This allows for containers to be made that hold all necessary software for running the application. Utilizing this technology allows for hosting on any operating system.
* We are also utilizing bootstrap which will allow for the front end to be used on different browsers

---

### Missing specifications

* Deleted Groups(Saved permissions)
* Edit Groups(Saved Permissions)
* Add Users to Groups
