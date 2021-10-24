# Serverless & Golang & Lambda

It is just a simple project written in golang for qualifying my ability to learn and write aws lambda functions with test cases.

## Installation

Use will need [serverless cli](https://www.serverless.com/framework/docs/providers/aws/guide/installation) to deploy lambda functions.



```bash
npm install -g serverless
```

You will need golang compiler

```bash
https://golang.org/doc/install
```



## Usage

If you want to build && test and deploy all functions use this command
```bash
makefile deploy
```


## Test Cases
Use this command to test written test cases in every lambda function folder.
```bash
go test -v .\devices\InsertDevice\InsertDevice.go .\devices\InsertDevice\InsertDevice_test.go
```


## License
[MIT](https://choosealicense.com/licenses/mit/)
