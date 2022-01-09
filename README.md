# Kitchen
Kitchen API which takes care of maximally efficient delivery of requested orders

## Setup & Run:
### Standard:
- change directory to root of Kitchen service
- run "go run src/main.go"
### Docker:
- change directory to root of Kitchen service
- run "docker build -t kitchen ."
- run "docker run -it -p 4006:4006 kitchen"
- check out port 4006