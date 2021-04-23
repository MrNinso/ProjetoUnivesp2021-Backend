docker rmi -f projeto-univesp2021-web:latest
docker rm -f projeto-univesp2021-db
docker rm -f projeto-univesp2021-web
docker run -d -p 3306:3306 --name projeto-univesp2021-db -e MYSQL_ROOT_PASSWORD=123 mariadb:latest
docker cp Database.sql projeto-univesp2021-db:/
docker build -t projeto-univesp2021-web:latest .
docker exec projeto-univesp2021-db /bin/sh -c "mysql -u root -p123 < /Database.sql"
docker run -it -e DATABASE_HOST="172.17.0.1" -e DATABASE_PORT="3306" -e DATABASE_USERNAME="root" -e DATABASE_PASSWORD="123" -e DATABASE_NAME="ProjetoUnivesp2021" -p 8080:8080 --name projeto-univesp2021-web projeto-univesp2021-web:latest