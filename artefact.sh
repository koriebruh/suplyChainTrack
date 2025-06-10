go get -u github.com/gofiber/fiber/v2
go get github.com/stretchr/testify/assert
go get github.com/json-iterator/go
go get github.com/gofiber/fiber/v2/middleware/cors
go get github.com/joho/godotenv
go get github.com/gofiber/fiber/v2/middleware/recover
go get github.com/gofiber/fiber/v2/middleware/helmet
go github.com/gofiber/fiber/v2/middleware/recover


#go install -tags 'mysql,postgres,mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

#create migrate
migrate create -ext sql -dir db/migrations create_table_first
migrate -database "mysql://root:korie123@tcp(mysql-db:3306)/suply_chain_track" -path db/migrations up


# go documentation swagegr
go get -u github.com/swaggo/files
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/gin-swagger/swaggerFiles

swag init -g cmd/main.go
tambahkan annotation ke handler lalu swag init lagi unutk update dokumentasi

localhost:8080/docs/index.html # untuk cek hasil dokumnetasi

docker compose -f docker/docker-compose.yml config --services

# Prometheus implementation go
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp