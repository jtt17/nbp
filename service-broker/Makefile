# Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

IMAGE = opensdsio/service-broker
VERSION := latest

all:build
.PHONY: all

build:service-broker
.PHONY: build

package:
	go get -v github.com/opensds/nbp/service-broker
.PHONY: package

service-broker:package
	mkdir -p  ./.output/
	go build -a -v -o ./.output/service-broker github.com/opensds/nbp/service-broker
.PHONY: service-broker

docker:build
	cp ./.output/service-broker .
	docker build . -t $(IMAGE):$(VERSION)
.PHONY: docker

clean:
	rm -rf ./.output/* ./service-broker
.PHONY: clean
