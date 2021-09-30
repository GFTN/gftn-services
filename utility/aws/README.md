# aws utility functions

util functions to aws parameter-store & secret-manager


# Configuration
---
## Environment Variables
Three new environment variables added to make AWS SDK work:
`AWS_ACCESS_KEY_ID` & `AWS_SECRET_ACCESS_KEY` & `AWS_REGION`

### how to get these values?
for `AWS_ACCESS_KEY_ID` & `AWS_SECRET_ACCESS_KEY` :
go to aws worldwire console -> click `services` -> click `IAM` -> click `Users` -> click your username or add a new user -> click `security credentials` tag -> click `create access key` -> export them as your env variables!

for `AWS_REGION`:
visit: https://github.com/jsonmaur/aws-regions, and the `Region Code` will be the expected input
# Implementation
---
## for `golang` and `javascript` users:

you'll need to define this struct so that the util function can identify which parameter/secret you want to access, this will be re-format as `/dev/p1.worldwire.io/crypto-service/secret`, and the `AWS_REGION` env variable is used to grab secret from the certain region, so make sure to also initialize the region env before calling these functions

example in `golang`
```
var test = utility.CredentialInfo{
	Environment: "dev",
	Domain:      "p1.worldwire.io",
	Service:     "crypto-service",
	Variable:    "secret",
}
```
then these are the secret content & parameter content for creating/updating secret/parameter
```
var newParameterContent = utility.ParameterContent{
	Value:       "newValue",
	Description: "yoyoyo",
}
var newTestContent = utility.SecretContent{
	Key:         "newKey",
	Value:       "newValue",
	Description: "yoyoyo again",
}
```
note that secret content can also be read from a file! here's how
```
var content = utility.SecretContent{
	FilePath:    "./test.json",
	Description: "this is a file, that's all",
}
```


example in `nodejs(typescript)`:
```
  let title: Var.CredentialInfo = {
    environment: "dev",
    domain: "worldwire.io",
    service: "IBM",
    variable: "test8",
  }

  let content: Var.ParameterContent = {
    value: "worldwire.io's parameter",
    description: "this is my parameter!"
  }

  let newContent: Var.SecretContent = {
    key: "this is",
    value: "new!",
    //filePath: "/Users/your.user/go/src/github.ibm.com/gftn/world-wire-services/utility/aws/javascript/src/test.json",
    description: "IBM"
  }
```
note that the filepath in `javascript` version needs to be absolute path!
and if you choose filepath, then `key` & `value` declaration won't be needed


## for `javascript` users:
go to `gftn-services/utility/aws/javascript` then run `npm install && npm run build` before using it
then import these lib functions
```
import * as PS from '<work_dir>/github.com/GFTN/gftn-services/utility/aws/javascript/build/awsParameter'
import * as PS from '<work_dir>/github.com/GFTN/gftn-services/utility/aws/javascript/build/awsSecret'
import * as Var from '<work_dir>/github.com/GFTN/gftn-services/utility/aws/javascript/build/utility/var'
```
and each utility function will return `promise`, which means that you can have either one of the following call to handle the asynchronous result
first one:
```
  let title: Var.CredentialInfo = {
    environment: "dev",
    domain: "worldwire.io",
    service: "IBM",
    variable: "test8",
  }

  let newContent: Var.ParameterContent = {
    value: "newest worldwire.io",
    description: "IBM New"
  }

  PS.createParameter(title, newContent)
  .then((res)=>{
    console.log(res)
    return PS.getParameter(title)
  })
  .then((res)=>{
    console.log(res)
    return PS.removeParameter(title)
  })
  .catch((err)=>{
    console.log(err)
  })
```
or this one using `async/await`
```
  let result: any
  try{
    result = await PS.getParameter(title)
  }catch(e){
    console.log(e)
  }
```


**note for whoever wants to create a new IAM user!**
To create a new IAM user, you will need to add the following permission to the users to access the utility functions:
```
IAM user required permission to call GetSecret function:
	* secretsmanager:GetSecretValue

	* kms:Decrypt - required only if you use a customer-managed AWS KMS key
	to encrypt the secret. You do not need this permission to use the account's
	default AWS managed CMK for Secrets Manager.

IAM user required permission to call UpdateSecret function:
	* secretsmanager:UpdateSecret

	* kms:GenerateDataKey - needed only if you use a custom AWS KMS key to
	encrypt the secret. You do not need this permission to use the account's
	AWS managed CMK for Secrets Manager.

	* kms:Decrypt - needed only if you use a custom AWS KMS key to encrypt
	the secret. You do not need this permission to use the account's AWS managed
	CMK for Secrets Manager.

IAM user required permission to call CreateSecret function:
	* secretsmanager:CreateSecret

	* kms:GenerateDataKey - needed only if you use a customer-managed AWS
	KMS key to encrypt the secret. You do not need this permission to use
	the account's default AWS managed CMK for Secrets Manager.

	* kms:Decrypt - needed only if you use a customer-managed AWS KMS key
	to encrypt the secret. You do not need this permission to use the account's
	default AWS managed CMK for Secrets Manager.

	* secretsmanager:TagResource - needed only if you include the Tags parameter.

IAM user required permission to call DeleteSecret function:
	* secretsmanager:DeleteSecret
	
IAM user required permission to call AppendSecret function:
	* secretsmanager:PutSecretValue

	* kms:GenerateDataKey - needed only if you use a customer-managed AWS
	KMS key to encrypt the secret. You do not need this permission to use
	the account's default AWS managed CMK for Secrets Manager.	
	
Note:
	recoveryDays should be 7 days at minimum
```
