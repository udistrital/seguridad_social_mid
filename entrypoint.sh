#!/usr/bin/env bash

set -e
set -u
set -o pipefail

if [ -n "${PARAMETER_STORE:-}" ]; then
  #No se conecta a base de datos
  #export SS_MID_API__PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}/ss_mid_api/db/username --output text --query Parameter.Value)"
  #export SS_MID_API__PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/ss_mid_api/db/password --output text --query Parameter.Value)"
fi

exec ./main "$@"
