# Install Air to Go path
# > curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# or install to ./bin/
# curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
# or install with directly with Go
# > go install github.com/cosmtrek/air@latest

# Air usage with Docker:
# docker run -it --rm \
#     -w "<PROJECT>" \
#     -e "air_wd=<PROJECT>" \
#     -v $(pwd):<PROJECT> \
#     -p <PORT>:<APP SERVER PORT> \
#     cosmtrek/air
#     -c <CONF>

dev:
	docker build --target=builder -t didstopia/nuttyqt:development .
	docker run --rm -it -v $(PWD):/app:delegated -w /app didstopia/nuttyqt:development

prod:
	docker build -t didstopia/nuttyqt:latest .
	docker run --rm -it didstopia/nuttyqt:latest
