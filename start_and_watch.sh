#!/bin/bash

docker compose up -d

docker wait exbestfriend

docker compose down