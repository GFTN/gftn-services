checkerrors()
{

    errors=$?

    # helps to alert failure if code does not build properly
    # and helps vscode debuger exit on typscript/lint error
    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: build process error, see deployment_prelaunch.sh"
        exit 1
    fi

}

# transpile helper
tsc ./project-provisioning/src/copy_file.ts

# copy dependent typescript files project:

comment='// ' \
fromPath='./authentication/src/models/auth.model.ts' \
toPath='./deployment/src/shared/models/auth.model.ts' \
searchTxt='./shared/models' \
replaceTxt='' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='./authentication/src/environment.ts' \
toPath='./deployment/src/shared/environment.ts' \
searchTxt='/shared/' \
replaceTxt='/' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/token.interface.d.ts' \
toPath='./deployment/src/shared/models/token.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/user.interface.d.ts' \
toPath='./deployment/src/shared/models/user.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/account.interface.d.ts' \
toPath='./deployment/src/shared/models/account.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/approval.interface.d.ts' \
toPath='./deployment/src/shared/models/approval.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/participant.interface.d.ts' \
toPath='./deployment/src/shared/models/participant.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/asset.interface.d.ts' \
toPath='./deployment/src/shared/models/asset.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/node.interface.d.ts' \
toPath='./deployment/src/shared/models/node.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='./secret-mgmt/src/encrypt.ts' \
toPath='./deployment/src/shared/encrypt.ts' \
node ./project-provisioning/src/copy_file.js

#transpile typescript
tsc -p ./deployment/tsconfig.deployment.debug.json

checkerrors

printf $(tput setaf 4)"\nTranspiled ./deployment \n\n"$(tput sgr0)
