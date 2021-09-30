checkerrors()
{

    errors=$?

    # helps to alert failure if code does not build properly
    # and helps vscode debuger exit on typscript/lint error
    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: build process error, see prelaunch.sh"
        exit 1
    fi

}

# traspile helper
tsc ./project-provisioning/src/copy_file.ts

# copy dependent typescript files project:

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/token.interface.d.ts' \
toPath='./authentication/src/shared/models/token.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/user.interface.d.ts' \
toPath='./authentication/src/shared/models/user.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/account.interface.d.ts' \
toPath='./authentication/src/shared/models/account.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/approval.interface.d.ts' \
toPath='./authentication/src/shared/models/approval.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/participant.interface.d.ts' \
toPath='./authentication/src/shared/models/participant.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/asset.interface.d.ts' \
toPath='./authentication/src/shared/models/asset.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='../gftn-web/src/app/shared/models/node.interface.d.ts' \
toPath='./authentication/src/shared/models/node.interface.d.ts' \
node ./project-provisioning/src/copy_file.js

comment='// ' \
fromPath='./secret-mgmt/src/encrypt.ts' \
toPath='./authentication/src/shared/encrypt.ts' \
node ./project-provisioning/src/copy_file.js

# gen-singleton-swagger & # gen-singleton-routes
# IMPORTANT: tsoa must run after copy_files.ts above
tsoa swagger -c ./authentication/tsoa.json
tsoa routes -c ./authentication/tsoa.json

#transpile typescript
BUILD="with js maps for local debug"
if [[ $* == *--prod* ]]
then
    tsc -p ./authentication/tsconfig.authentication.prod.json
    BUILD="for cloud deployment"
else
    tsc -p ./authentication/tsconfig.authentication.debug.json
fi

checkerrors

printf $(tput setaf 4)"\nTranspiled ./authentication ${BUILD} \n\n"$(tput sgr0)
