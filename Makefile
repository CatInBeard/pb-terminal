APP_NAME = terminal

all: build install

build: # Build application with docker
	docker run --rm -v $(PWD):/app 5keeve/pocketbook-go-sdk build -o $(APP_NAME).app
install: # Copy the application to the cable-connected PB
	mv $(APP_NAME).app $(shell find /media /mnt /run/media -type d -name '*PB*' -print -quit 2>/dev/null)/applications/$(APP_NAME).app && echo -e "\033[32mThe application is successfully installed on your PB device\033[0m" || echo -e "\033[31mFailed to copy application on your PB device, check connection\033[0m"

clean:
	rm $(APP_NAME).app