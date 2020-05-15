## Cells

Cells are the equivalent of buildpacks. They create "cells" for services that encapsulate 
their dependencies and runtime requirement. They isolate the service from the outside 
world and provide a single entry point via http port 8080.

In the event a service does not have a http server e.g shell scripts we start one for it.
