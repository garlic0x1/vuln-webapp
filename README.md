# vuln-webapp
Vulnerable webapp for testing  
There are 2 XSS, a SQL injection, a local file inclusion, and broken access control (that I know of)  

# setup
```
docker-compose build
docker-compose up
```
This starts a mysql db at port 3306 and the webapp at port 9999
