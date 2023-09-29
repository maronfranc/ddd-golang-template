# ddd-golang-template


### Request tests
#### Create example
```sh
DEV_API_HOST=http://localhost:3000
PROD_API_HOST=http://localhost:4001
API_HOST=${PROD_API_HOST}


curl localhost:3000/examples
curl localhost:3000/examples/123
curl -XPOST localhost:3000/examples \
  -H "Content-Type: application/json" \
  --data '{"title": "Manual_test_example", "description": "Manual description test"}'
```
