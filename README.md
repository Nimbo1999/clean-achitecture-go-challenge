## Example command for running evans with docker

docker run --rm -it -v "$(pwd):/mount:ro" ghcr.io/ktr0731/evans:latest -r --host 192.168.100.14
