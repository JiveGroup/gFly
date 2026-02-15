module gfly

go 1.24.0

require (
	github.com/gflydev/cache v1.0.5
	github.com/gflydev/console v1.1.0
	github.com/gflydev/core v1.18.0
	github.com/gflydev/db v1.14.0
	github.com/gflydev/db/psql v1.4.9
	github.com/gflydev/http v1.0.2
	github.com/gflydev/modules/storage v1.0.6
	github.com/gflydev/modules/storagecs3 v1.1.1
	github.com/gflydev/notification v1.1.1
	github.com/gflydev/notification/mail v1.0.3
	github.com/gflydev/session v1.0.3
	github.com/gflydev/session/redis v1.0.3
	github.com/gflydev/storage v1.1.6
	github.com/gflydev/storage/local v1.1.7
	github.com/gflydev/utils v1.1.0
	github.com/gflydev/view/pongo v1.0.3
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/swaggo/swag v1.16.6
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/flosch/pongo2/v6 v6.0.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/gflydev/mail v1.0.3 // indirect
	github.com/gflydev/storage/cs3 v1.2.3 // indirect
	github.com/gflydev/validation v1.2.1 // indirect
	github.com/go-ini/ini v1.67.1 // indirect
	github.com/go-openapi/jsonpointer v0.22.4 // indirect
	github.com/go-openapi/jsonreference v0.21.4 // indirect
	github.com/go-openapi/spec v0.22.3 // indirect
	github.com/go-openapi/swag/conv v0.25.4 // indirect
	github.com/go-openapi/swag/jsonname v0.25.4 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.4 // indirect
	github.com/go-openapi/swag/loading v0.25.4 // indirect
	github.com/go-openapi/swag/stringutils v0.25.4 // indirect
	github.com/go-openapi/swag/typeutils v0.25.4 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hibiken/asynq v0.26.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jivegroup/fluentsql v1.5.4 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.98 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/redis/go-redis/v9 v9.17.3 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/tinylib/msgp v1.6.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.69.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	golang.org/x/tools v0.42.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// Replace deprecated module paths with correct ones
replace github.com/go-ini/ini => gopkg.in/ini.v1 v1.67.1
