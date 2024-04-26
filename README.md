#  CTI - Golang Base  

## I. Architecture [Xem tại đây](https://app.diagrams.net/#G1Gco0oD-ePPhw7MExIP_MpXHCI5XTk3TC)

### 1. cmd
#### 1.1 cmd/main.go
```
app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "configPrefix",
			Aliases:     []string{"confPrefix"},
			Usage:       "prefix for config",
			Value:       "auth",
			Destination: &configPrefix,
		},
		&cli.StringFlag{
			Name:        "configSource",
			Aliases:     []string{"confSource"},
			Value:       "../config/.env",
			Usage:       "set path to environment file",
			Destination: &configSource,
		},
	}
```
- Phần code trên define đường dẫn đến file env và prefix cho các biến trong env
- Các biến env được định nghĩa trong config/config.go
```
AUTH_LOGLEVEL=-1
AUTH_SERVICELOGFILE=./service.log
AUTH_DATABASELOGFILE=./database.log

AUTH_SERVICENAME=CTIGroupDev
AUTH_JWTSECRET=
```

### 2. src
#### 2.1 src/const 
- Các file trong const define các hằng số  global
#### 2.2 src/server
- Khởi tạo server graph và grpc
- Define router graph
```
func v1(r chi.Router) {
	configAdmin := generated_admin.Config{Resolvers: &resolver_admin.Resolver{}}
	configAdmin.Directives = directive.AdminDirective

	configUser := generated_user.Config{Resolvers: &resolver_user.Resolver{}}
	configUser.Directives = directive.UserDirective

	srvAdmin := handler.NewDefaultServer(generated_admin.NewExecutableSchema(configAdmin))
	srvUser := handler.NewDefaultServer(generated_user.NewExecutableSchema(configUser))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AllowAll().Handler)
		r.With(middleware.Middleware()).Route("/graphql", func(r chi.Router) {
			r.Handle("/admin", srvAdmin)
			r.Handle("/user", srvUser)
		})
	})
}
```
#### 2.3 src/middleware
#### 2.4 src/network
- Define một số hàm được sử dụng chung trong toàn bộ hệ thống 
#### 2.5 src/redis
- Dùng để connect và sử dụng redis 
#### 2.6 src/schedule 
- Xử lý một số cron job
#### 2.7 src/database
- File connect.go xử lý việc kết nối với database
```
    for i := 1; i <= config.Get().NumberRetry; i++ {
        mongoClient, err = mongo.NewMongoDBFromUrl(ctx, config.Get().MongoURL, time.Second*10)
        if err != nil {
            if i == config.Get().NumberRetry {
                return err
            }
            time.Sleep(10 * time.Second)
        }

        if mongoClient != nil {
            break
        }
    }
    
    // load account collection 
	if err := collection.LoadAccountCollectionMongo(mongoClient); err != nil {
		return err
	}
```
- Folder model define các bảng trong database 
- Collection xử lý việc kết nối đến các bảng trong database, đánh index cho các bảng đó

#### 2.8 src/graph
##### 2.8.1 schema
- Các model graph được define trong graph/model, gồm có 2 loại type (các model trả về khi call API) được define trong các file có đuôi .type.graphql và input (data được client truyền lên khi call API) được define trong các file có đuôi .input.graphql
- File general.graphql define các loại data chung như pagination hay các directive check auth
- 2 folder admin và user define các api được call với từng loại đối tượng là admin hay user

##### 2.8.2 generated
- Các API, model được define trong schema sẽ được generate trong folder này để có thể sử dụng dưới dạng graph
##### 2.8.3 resolver 
- Các API được define trong schema sẽ được xử lý tại đây
- Các input được khởi tạo theo command được define trong src/service, và được valid và xử lý trong src/controller
##### 2.8.4 directive
- Các directive trong schema sẽ được define cụ thể tại đây
```
RequiredAuthAdmin: func(ctx context.Context, obj interface{}, next graphql.Resolver, action []*string, actionAdmin *string, check_ip *bool) (res interface{}, err error) {
		if !network.HasToken(ctx) {
			return nil, fmt.Errorf("unauthorized")
		}

		tokenStr := network.Token(ctx)
		result, err := service_shared.TokenVerify(ctx, &service_shared.TokenVerifyCommand{Token: tokenStr})
		if err != nil || result == nil {
			return nil, err
		}

		if result.AccountType != AccountTypeAdmin {
			return nil, fmt.Errorf("permission deny")
		}

		ctx = context.WithValue(ctx, "workspace_id", result.WorkspaceID)
		ctx = context.WithValue(ctx, "account_id", result.AccountID)
		ctx = context.WithValue(ctx, "user_id", result.UserID)
		ctx = context.WithValue(ctx, "email", result.Email)
		ctx = context.WithValue(ctx, "sub_workspace_id", result.SubWorkspaceID)

		return next(ctx)
	},
```

#### 2.9 src/grpc
##### 2.9.1 protoc : là plugin dùng để gen code từ proto 
##### 2.9.2 proto : define các package, hàm, message 
##### 2.9.3 golang : chứa code được gen từ folder proto 
##### 2.9.4 grpc_client : define các func để connect đến server grpc của các service khác

## II. Triển khai
### 1. Quy trình code một API graph
- Định nghĩa type, input graphql trong src/graph/schema/model
- Định nghĩa api trong src/graph/schema/admin (đối với các API dành cho admin) và /user (đối với các API dành cho user)
- Định nghĩa model tại src/database/model
- Định nghĩa collection tại src/database/collection và load collection tại file connection.go 
- Xử lý logic tại service
- Xử lý tại controller
- Xử lý tại resolver
  
### 2. Quy trình code một hàm grpc
- Copy folder protoc, file gen từ một service khác và paste vào folder grpc
- Tạo folder proto và define các func, message
- Chạy lệnh gen
- Tạo file với cú pháp tên service_grpc_server.go để define grpc server cho service(có rồi thì bỏ qua bước này)
- Logic của các hàm này được xử lý tại service/shared/grpc (nếu chưa có thì tạo theo path trên)
- Define các hàm có tên trùng khớp với rpc đã define ở proto và gọi các hàm ở service shared để xử lý và trả về
- Tạo folder grpc_client(nếu đã có thì bỏ qua bước này)
- Trong grpc_client define các hàm để connect với server grpc của service khác
#### Lưu ý : với struct hiện tại, các rpc và message phải được triển khai ở cả 2 service