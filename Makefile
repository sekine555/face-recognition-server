RM = rm -rf
DIR = ./docker/db/data/

# windowsで実行する場合
ifeq ($(OS),Windows_NT)
    RM = cmd.exe /C rd /S /Q
    DIR = .\docker\db\data\

    endif

up:
	docker-compose down --rmi all
	-$(RM) $(DIR)
	docker-compose up -d --build
