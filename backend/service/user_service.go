package service

import (
    "context"
    "database/sql"
    "errors"
    "time"

    pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb"
    _ "github.com/lib/pq" // PostgreSQL driver
    "golang.org/x/crypto/bcrypt"
    "google.golang.org/protobuf/types/known/timestamppb"
    "github.com/golang-jwt/jwt/v5"
)

type UserServiceServer struct {
    pb.UnimplementedUserServiceServer
    db *sql.DB  // PostgreSQL-yhteys
}


var jwtKey = []byte("salainen-avain") // vaihda turvalliseen

// Helper: Tarkistaa, että token vastaa käyttäjän ID:tä
func getUserIDFromContext(ctx context.Context) (int32, error) {
    userID, ok := ctx.Value("user_id").(int32)
    if !ok {
        return 0, errors.New("unauthorized")
    }
    return userID, nil
}

// Luo käyttäjän
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    if req.Name == "" || req.Email == "" || req.Password == "" {
        return nil, errors.New("name, email and password are required")
    }

    // Tarkista, ettei käyttäjä ole jo olemassa
    var exists bool
    err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("user with this email already exists")
    }

    // Hashaa salasana
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Lisää tietokantaan
    var id int32
    var createdAt time.Time
    err = s.db.QueryRow(
        "INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id, created_at",
        req.Name, req.Email, string(hashedPassword),
    ).Scan(&id, &createdAt)
    if err != nil {
        return nil, err
    }

    user := &pb.User{
        Id:        id,
        Name:      req.Name,
        Email:     req.Email,
        CreatedAt: timestamppb.New(createdAt),
    }

    return &pb.CreateUserResponse{User: user}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    userID, err := getUserIDFromContext(ctx)
    if err != nil {
        return nil, err
    }

    // Varmistetaan, että käyttäjä saa vain omat tietonsa
    if req.Id != userID {
        return nil, errors.New("cannot access another user's data")
    }

    var user pb.User
    var createdAt time.Time
    err = s.db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id=$1", userID).
        Scan(&user.Id, &user.Name, &user.Email, &createdAt)
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    } else if err != nil {
        return nil, err
    }

    user.CreatedAt = timestamppb.New(createdAt)
    return &user, nil
}

// Kirjautuminen
func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    if req.Email == "" || req.Password == "" {
        return nil, errors.New("email and password are required")
    }

    var user pb.User
    var hashedPassword string
    var createdAt time.Time

    err := s.db.QueryRow(
        "SELECT id, name, email, password, created_at FROM users WHERE email=$1",
        req.Email,
    ).Scan(&user.Id, &user.Name, &user.Email, &hashedPassword, &createdAt)
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    } else if err != nil {
        return nil, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid password")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.Id,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return nil, err
    }

    user.CreatedAt = timestamppb.New(createdAt)
    return &pb.LoginResponse{
        Token: tokenString,
        User:  &user,
    }, nil
}

// Päivittää käyttäjän tiedot
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
    userID, err := getUserIDFromContext(ctx)
    if err != nil {
        return nil, err
    }
    if req.Id != userID {
        return nil, errors.New("cannot update another user's data")
    }

    if req.Name != "" {
        _, err = s.db.Exec("UPDATE users SET name=$1 WHERE id=$2", req.Name, userID)
        if err != nil {
            return nil, err
        }
    }
    if req.Email != "" {
        _, err = s.db.Exec("UPDATE users SET email=$1 WHERE id=$2", req.Email, userID)
        if err != nil {
            return nil, err
        }
    }

    var user pb.User
    var createdAt time.Time
    err = s.db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id=$1", userID).
        Scan(&user.Id, &user.Name, &user.Email, &createdAt)
    if err != nil {
        return nil, err
    }

    user.CreatedAt = timestamppb.New(createdAt)
    return &user, nil
}

// Poistaa käyttäjän tilin
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    userID, err := getUserIDFromContext(ctx)
    if err != nil {
        return nil, err
    }
    if req.Id != userID {
        return nil, errors.New("cannot delete another user's account")
    }

    res, err := s.db.Exec("DELETE FROM users WHERE id=$1", userID)
    if err != nil {
        return nil, err
    }
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return nil, errors.New("user not found")
    }

    return &pb.DeleteUserResponse{Success: true}, nil
}
