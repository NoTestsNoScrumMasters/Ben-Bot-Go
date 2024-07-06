container_name := Ben-Bot-Go

build:
	docker build -t [container_name]

run:
	docker run --name ${container_name} -p 3000:3000 ${container_name}

start:
	docker run -d --name ${container_name} -p 3000:3000 ${container_name}

stop:
	docker container stop ${container_name}