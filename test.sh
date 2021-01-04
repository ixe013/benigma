#!/bin/bash
vault list benigma/models
vault write benigma/models name=TEST_MODEL
vault list benigma/models
vault delete benigma/models/TEST_MODEL

vault list benigma/instances
vault write benigma/models/M3/instance id=first
vault list benigma/instances
vault read benigma/instances/first
