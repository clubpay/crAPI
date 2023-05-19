#! /bin/sh

#
# Licensed under the Apache License, Version 2.0 (the “License”);
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an “AS IS” BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


## Uncomment the following line if you wish to run tests
IS_TESTING=True python3 manage.py test --no-input &&\
python3 manage.py migrate user --fake &&\
python3 manage.py migrate crapi &&\
python3 manage.py migrate db

if [ "$TLS" != "" ]; then
  # if $TLS_CERTIFICATE and $TLS_KEY are not set, use the default ones
  if [ "$TLS_CERTIFICATE" == "" ]; then
    TLS_CERTIFICATE=$PWD/certs/server.crt
  fi
  if [ "$TLS_KEY" == "" ]; then
    TLS_KEY=$PWD/certs/server.key
  fi
  python3 manage.py runsslserver 0.0.0.0:${SERVER_PORT} --certificate "$TLS_CERTIFICATE" --key "$TLS_KEY"
else
  python3 manage.py runserver 0.0.0.0:${SERVER_PORT}
fi
exec "$@"
