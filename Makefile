# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

APPNAME = cloud-logging-log-creator
PROJECT = $(shell gcloud config get-value project)

env:
	gcloud config set project $(PROJECT)

image:
	docker build -t gcr.io/$(PROJECT)/$(APPNAME) .

imagear:
	docker build -t us-docker.pkg.dev/$(PROJECT)/$(APPNAME)/latest .	

serve: 
	docker run --name=$(APPNAME) -d -P $(APPNAME)	

clean:
	-docker stop $(APPNAME)
	-docker rm $(APPNAME)
	-docker rmi $(APPNAME)		

push:
	docker push gcr.io/$(PROJECT)/$(APPNAME)	

pushar: 
	docker push us-docker.pkg.dev/$(PROJECT)/$(APPNAME)/latest

publish: clean image push 	