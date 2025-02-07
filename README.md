how to run:

```
go run cmd/main.go
```

using hot reload golang with .air.toml

```
air
```

using hot reload without .air.toml

```
- air -c ./cmd/main.go
```

to build frontend:

```
cd web
npm run build
```

generate docs swagger:

```
swag init -g cmd/main.go
```
