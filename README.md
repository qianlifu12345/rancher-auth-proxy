rancher-auth-filter-service
========

A microservice for validate the token and return account id and environment id.

## Building

`make`

##Parameter
`   --rancherUrl value  Rancher server url (default: "http://54.255.182.226:8080/") [$RANCHER_SERVER_URL]`
`   --localport value   Local server port  (default: "8092") [$LOCAL_VALIDATION_FILTER_PORT]`
   
## Running

`./bin/rancher-auth-filter-service    --rancherUrl <value>    --localport <value> `

Or set the environment variables
`[$RANCHER_SERVER_URL]`
`[$LOCAL_VALIDATION_FILTER_PORT]`

## License
Copyright (c) 2014-2016 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
# rancher-auth-proxy
