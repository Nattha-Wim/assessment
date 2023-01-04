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

* Postman
  * run  DATABASE_URL="postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb" go run server.go
  * use expenses.postman_collection_env.json for run test


## Unsuccessful
* docker-compose testing sandbox (integration test)
    * problem : 
      Can't create http://localgost:2565/expenses because connection refused. I try to slove it but it not success and I guess the problem is container and network are working separately.

      

