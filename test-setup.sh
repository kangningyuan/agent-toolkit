#!/bin/bash

echo "Testing Docker setup..."

# 测试镜像是否存在
echo "=== Checking Docker images ==="
docker images | grep -E "(alpine|python)"

# 测试容器运行
echo ""
echo "=== Testing Alpine container ==="
docker run --rm swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/alpine:latest echo "Alpine test successful"

echo ""
echo "=== Testing Python container ==="
docker run --rm swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/python:3-alpine python -c "print('Python test successful')"

echo ""
echo "=== Setup test complete ==="
