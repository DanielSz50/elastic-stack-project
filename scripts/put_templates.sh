#!/bin/bash

curl -X PUT "localhost:9200/_component_template/app_component_template?pretty" -H "Content-Type: application/json" -d "$(cat template/app-component-template.json)"
curl -X PUT "localhost:9200/_index_template/app_index_template?pretty" -H "Content-Type: application/json" -d "$(cat template/app-index-template.json)"
