module github.com/example/git-ops/backend

go 1.23

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/viper v1.18.2
	github.com/valyala/fasthttp v1.51.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)

require github.com/gofiber/fiber/v2 v2.52.5
