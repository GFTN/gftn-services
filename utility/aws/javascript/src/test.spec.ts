// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as PS from './awsParameter';
import * as SM from './awsSecret';
import * as Var from './utility/var';
import * as Env from './env';

env_main();


function env_main(){
  process.env['HOME_DOMAIN_NAME'] = "p1.worldwire.io";
  process.env['SERVICE_NAME'] = "api-service";
  process.env['ENVIRONMENT_VERSION'] = "dev";
  process.env['ACCOUNT_STORAGE_LOCATION'] = "dev";
  process.env['ACCOUNT_SOURCE'] = "dev";
  Env.InitEnv();
}

function parameter_main() {
  
  const title: Var.CredentialInfo = {
    environment: "dev",
    domain: "worldwire.io",
    service: "IBM",
    variable: "test9"
  };

  const oldContent: Var.ParameterContent = {
    value: "worldwire.io",
    description: "IBM"
  };


  /* await/async way(ensure that you add the async keyword at this function)
  let result: any
  try{
    result = await PS.getParameter(title)
  }catch(e){
    console.log(e)
  }
    console.log(result)
    console.log("hihih")
  */

  const newContent: Var.ParameterContent = {
    value: "newest worldwire.io",
    description: "IBM New"
  };

  PS.createParameter(title, oldContent)
    .then((res) => {
      console.log(res);
      return PS.getParameter(title);
    })
    .then((res) => {
      console.log(res);
      return PS.updateParameter(title, newContent);
    })
    .then((res) => {
      console.log(res);
      return PS.getParameter(title);
    })
    .then((res) => {
      console.log(res);
      return PS.removeParameter(title);
    })
    .then((res) => {
      console.log(res);
      return PS.getParameter(title);
    })
    .catch((err) => {
      console.log(err);
    });

}

function secret_main() {
  
  const title: Var.CredentialInfo = {
    environment: "dev",
    domain: "worldwire.io",
    service: "IBM",
    variable: "test12"
  };
  
  const oldContent: Var.SecretContent = {
    // key: "dev",
    // value: "worldwire.io",
    filePath: "/Users/your.user/go/src/github.ibm.com/gftn/world-wire-services/utility/aws/javascript/src/test.json",
    description: "IBM"
  };

  const newContent: Var.SecretContent = {
    key: "this is",
    value: "new!",
    //filePath: "/Users/your.user/go/src/github.ibm.com/gftn/world-wire-services/utility/aws/javascript/src/test.json",
    description: "IBM"
  };


  SM.createSecret(title, oldContent)
    .then((res) => {
      console.log(res);
      return SM.getSecret(title);
    })
    .then((res) => {
      console.log(res);
      return SM.updateSecret(title, newContent);
    })
    .then((res) => {
      console.log(res);
      return SM.getSecret(title);
    })
    .then((res) => {
      console.log(res);
      return SM.removeSecret(title);
    })
    .then((res) => {
      console.log(res);
      return SM.getSecret(title);
    })
    .catch((err) => {
      console.log(err);
    });


}
