docker rmi -f ProjetoUnivesp2021WEB:latest
docker rm -f ProjetoUnivesp2021DB
docker build -t ProjetoUnivesp2021WEB:latest .
docker run -d -p 3306:3306 --name ProjetoUnivesp2021DB -e MYSQL_ROOT_PASSWORD=123  mariadb:latest
docker cp Datbase.sql ProjetoUnivesp2021DB:/
docker exec ProjetoUnivesp2021DB /bin/sh -c "mysql --host=127.0.0.1 --port=3306 -u root -p123 < Database.sql"
docker run docker run --it -e DATABASE_HOST="172.17.0.1" -e DATABASE_PORT="3306" -e DATABASE_USERNAME="root" -e DATABASE_PASSWORD="123" -e DATABASE_NAME="ProjetoUnivesp2021" -p 8080:8080 --name ProjetoUnivesp2021-WEB deletar:latest