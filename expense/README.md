## Manual
* Authentication
  * User: admin
  * Password : 45678

* URL
  * GET http://localhost:2565/
    * home
  * POST http://localhost:2565/expenses 
    * for create expenses
  * GET http://localhost:2565/expenses
    * for get all expenses
  * GET http://localhost:2565/expenses/id
    * for get only one expenses by id
  * PUT http://localhost:2565/expenses/id
    * for update expenses by id
 
* Unit Test
  * go test -tags=unit -v ./...

* Integration Test 
  * DATABASE_URL="postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb" go run server.go
  * AUTH_TOKEN="_____" go test --tags=integration -v ./...   
  
* Multi-Stage build ( use below 3 command )
  * docker build -t devops:assessment .
  * docker run devops:assessment 
  * docker images devops:assessment 

* Docker-compose testing sandbox ( use below 2 command )
  * docker-compose -f docker-compose.test.yml down    
  * docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests
  * 
* Postman
  * run  DATABASE_URL="postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb" go run server.go
  * use expenses.postman_collection_env.json for run test




      

