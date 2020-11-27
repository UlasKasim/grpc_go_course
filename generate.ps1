#!/bin/bash

protoc  --go_out=plugins=grpc:. greet/greetpb/greet.proto ` calculator/calculatorpb/calculator.proto