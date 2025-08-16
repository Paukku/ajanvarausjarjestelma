package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "log"
    "net"
    "os"

    pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func init() {
    // Lataa .env tiedosto (vain kehityksessä)
    err := godotenv.Load("backend/.env")
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is not set")
    }
    jwtKey = []byte(secret)
}

// AuthInterceptor tarkistaa JWT-tunnuksen
func AuthInterceptor(ctx context.Context) (context.Context, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, errors.New("missing metadata")
    }
    tokenStr := md["authorization"]
    if len(tokenStr) == 0 {
        return nil, errors.New("missing token")
    }

    token, err := jwt.Parse(tokenStr[0], func(t *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }

    claims := token.Claims.(jwt.MapClaims)
    userID := int32(claims["user_id"].(float64))
    ctx = context.WithValue(ctx, "user_id", userID)
    return ctx, nil
}

// UserServiceServer on sama kuin aiemmin
type UserServiceServer struct {
    pb.UnimplementedUserServiceServer
    db *sql.DB
}

func main() {
    // Yhdistetään PostgreSQL:ään
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL is not set")
    }
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("failed to connect to DB: %v", err)
    }
    defer db.Close()

    // Testataan yhteys
    if err := db.Ping(); err != nil {
        log.Fatalf("DB ping failed: %v", err)
    }
    fmt.Println("Connected to PostgreSQL!")

    // Käynnistetään gRPC-palvelin
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(func(
            ctx context.Context,
            req interface{},
            info *grpc.UnaryServerInfo,
            handler grpc.UnaryHandler,
        ) (interface{}, error) {
            // Suojaa UpdateUser ja DeleteUser metodit JWT:llä
            if info.FullMethod == "/pb.UserService/UpdateUser" ||
                info.FullMethod == "/pb.UserService/DeleteUser" {
                ctx, err = AuthInterceptor(ctx)
                if err != nil {
                    return nil, err
                }
            }
            return handler(ctx, req)
        }),
    )

    // Rekisteröidään palvelin
    pb.RegisterUserServiceServer(grpcServer, &UserServiceServer{db: db})

    fmt.Println("gRPC server running on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
