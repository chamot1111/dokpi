0.1
---
* bootstrap.sh
* CLI:
  * build repo
  * deploy repo
  * stop repo
* repo framework: golang
* detect change of git repo, restart repo
* ENV variable: $DOKPI_DATAFOLDER, $DOKPI_PORT
* services command line output is logged

0.2
---
* web interface for services status

0.3
---
* hot backup through http: get a zip of the data folder without stopping the services

0.4
---
* cold backup through http:
  1- stop the services
  2- get a zip of the data folder
  3- start the services

0.5
---
* repo framework: golang, node.js

0.6
---
* web interface control services (start, stop)

0.7
---
* proxy with domain name
