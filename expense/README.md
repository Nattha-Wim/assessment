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
  
* Multi-Stage build ( use the 3 commands below here)
  * docker build -t devops:assessment .
  * docker run devops:assessment 
  * docker images devops:assessment 

* Docker-compose testing sandbox ( use the 2 commands below here)
  * docker-compose -f docker-compose.test.yml down    
  * docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests
  * 
* Postman
  * run  DATABASE_URL="postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb" go run server.go
  * use expenses.postman_collection_env.json for run test
![result](image/postman.png)
### Problem
    First time I separate braches to create each story. When it succeeds I merge it into main branch.
    I do this process for every story until complete evey story and main branch finally work. 
    So I imporove everything in main branch.
    But at the end, I realized that I should better improve the code in each branch not the main branch.

  ![result](image/git-log-graph.png)