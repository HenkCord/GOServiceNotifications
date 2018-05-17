Service Notifications Golang
===
example microservice

Start programm
---
`go run main.go`

#### Flags
--environment=**production**, **development**, **local** 
##### Example
`go run main.go --environment=production`

Start programm with watch
---
`nodemon ./main.go`

---
## Database
**notifications**
* emailTemplates
  * _id
  * name
  * description
  * subject 
    * ru
    * en 
    * ...
  * message
    * ru
    * en 
    * ...
  * updateAt
  * createdAt
* pushTemplates
  * _id
  * name
  * description
  * title 
    * ru
    * en 
    * ...
  * message
    * ru
    * en 
    * ...
  * updateAt
  * createdAt
* smsTemplates
  * _id
  * name
  * description
  * message
    * ru
    * en 
    * ...
  * updateAt
  * createdAt