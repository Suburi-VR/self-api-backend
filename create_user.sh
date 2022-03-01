#!/bin/bash -eu

aws --endpoint http://192.168.1.8:9229 cognito-idp admin-create-user --user-pool-id "local_5HgXw3xJ" --username $1 --user-attributes Name="custom:supporter",Value="true"

exit 0