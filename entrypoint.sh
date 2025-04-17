#!/bin/sh
set -e

echo "Running migrations..."
usdc-transfer-listener-svc migrate up

echo "Running service"
usdc-transfer-listener-svc run service