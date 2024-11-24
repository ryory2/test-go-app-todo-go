# 利用バージョン

### バージョン
golang:1.23.2
gin@v1.10.0
gorm@v1.25.1
postgres@v1.5.0
testify@v1.8.5
golangci-lint@v1.50.1
golang-migrate@4.15.2

### 出張方法
githubでプロジェクトを作成
※git initで.gitファイルが有ること
git bash
git remote add origin https://github.com/ryory2/test-go-app-todo-go.git
鍵を設定
tortoise git でプッシュ

go mod init github.com/ryory2/test-go-app-todo-go　（go.modファイルが作られる）
mkdir -p cmd/server
mkdir -p internal/handler
mkdir -p internal/repository
mkdir -p internal/service
mkdir -p internal/model
mkdir -p migrations
mkdir -p tests

#### 4. 必要なライブラリとツールの導入
go get -u github.com/gin-gonic/gin@v1.10.0
go get -u gorm.io/gorm@v1.25.1
go get -u gorm.io/driver/postgres@v1.5.0
go get -u github.com/stretchr/testify@v1.8.5
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
<!-- -tagsは利用するDBにより変更する -->
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
設定ファイル .golangci.yml をプロジェクトルートに追加
touch .golangci.yml
<!-- curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz -o migrate.tar.gz
tar -xzf migrate.tar.gz
sudo mv migrate /usr/local/bin/
migrate -version -->

## データベースの設定
mkdir -p model
touch model/task.go

### マイグレーションの作成
touch migrations/202411231200_create_tasks_table.up.sql
touch migrations/202411231200_create_tasks_table.down.sql

### マイグレーションの実行
migrate -path=./migrations -database "postgres://admin:password@localhost:5432/todo_db?sslmode=disable" up
（削除する場合）migrate -path=./migrations -database "postgres://admin:password@localhost:5432/todo_db?sslmode=disable" down

### APIエンドポイントの作成
touch internal/repository/repository.go
touch internal/service/task_service.go
touch internal/handler/task_handler.go
touch cmd/server/main.go
touch tests/task_service_test.go

### Lintの実行
golangci-lint run







## ユニットテスト
touch internal/repository/mock_repository.go
touch internal/handler/task_test.go