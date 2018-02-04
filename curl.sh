echo --- post to /login/
curl -H "Content-Type: application/json" -X POST -d '{"username":"admin","password":"pass123"}' http://localhost:8080/api/v1/login/
curl -H "Content-Type: application/json" -X POST -d '{"username":"joklad","password":"pass123"}' http://localhost:8080/api/v1/login/
curl -H "Content-Type: application/json" -X POST -d '{"username":"admin","password":"1234"}' http://localhost:8080/api/v1/login/


