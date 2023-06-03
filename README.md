# Create and Deploy AWS Lambda Function with Go1.X
Services - AWS Lambda, DynamoDB, Cloudwatch, APIGateway

---

References
- https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/index.html
- https://awscli.amazonaws.com/v2/documentation/api/latest/reference/iam/index.html

- **Make sure to create and configure AWS account in your local environment**
- **Make sure to replace the `<account_id>` in lambda-policy.json and role-policy.json files**
* Install Dependencies - `go get -v all`
* Build Command - `GOOS=linux go build -o build/main cmd/main.go`
* Zip Command - `zip -jrm build/main.zip build/main`

### Create Policy & Role and Attach

- Create role - `aws iam create-role --role-name go-lambda-executor --assume-role-policy-document file://role-policy.json`
- Update required policies - `aws iam put-role-policy --role-name go-lambda-executor --policy-name go-lambda-policy --policy-document file://lambda-policy.json`

### Create Functions

* Get the **Arn** for next command - `aws iam get-role --role-name go-lambda-executor`
* Create Function - `aws lambda create-function --function-name go-lambda-ex --runtime go1.x --handler main --zip-file fileb://build/main.zip --role <Arn>`
