Service mock
============

Service databases connections description.

## Run

## Docker compose

The creation of mocks can be started directly from the root directory of the project.  

    docker-compose -p goservice_mock -f mock/docker-compose.yml up -d

In case of conflicts in port numbers, they can be changed to others. In this case you need to remember to change them in the service and migrago configs.    

See that the required dependencies have started:  

    docker-compose -p goservice_mock -f mock/docker-compose.yml ps

Clean up after the completion of the project:

    docker-compose -p goservice_mock -f mock/docker-compose.yml down --remove-orphans

## Performing migrations

    migrago -c config_migrations-local.yaml up

## Запуск сервиса

    go run main.go -c config-local.yaml
