# World Wire Middleware
/authorization/**/* contains a custom built middelware package used by various IBM Blockchain World Wire micro-services. It is written in GOLANG and is only compatible with other golang micro-services. 

# Defined Permisisons
Required permissions are defined at `./permissions.yaml` (and transpiled into `./permissions.json` and consumed by the middleware)

# JWT vs. User Permission

## JWT 
Does not distinguish between super vs. particpant permissions. The server that possesses the JWT 
will have access to the requested resources if the JWT has sufficient claims. The permissions checking 
logic is implemented in the `middleware/token/token.go`.

## User Permission 
User Permissions are implmented using the firebase authentication service (which uses JWT under the hood). Permissions are grouped into either *super* or *participant*. The permission checking logic is 
implemented  in the `middleware/token/client-token.go`.

## Steps to run and test locally
The ./main.go is a demo api application that can be used for local development and testing. To debug locally:  
1. open vscode editor from `gftn-services/auth-service`
2. install dependencies:   
`$ cd ../ ; dep ensure`
NOTE: If this fails install dependencies one by one using: 
`$ go get "github.com/org/some_project_name"`  
You will have to manually check the debug output for which packages are missing.
3. Run the app in debug mode by selecting the play button in .vscode and select `launch authorization demo`