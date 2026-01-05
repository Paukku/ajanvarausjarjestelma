package validation

import (
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
)

var (
	emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

func ValidateCreateUserRequest(req *pb.CreateUserRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	// --- NAME ---
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if len(name) < 2 {
		return status.Error(codes.InvalidArgument, "name must be at least 2 characters")
	}
	if len(name) > 50 {
		return status.Error(codes.InvalidArgument, "name must be at most 50 characters")
	}

	// --- EMAIL ---
	email := strings.TrimSpace(strings.ToLower(req.GetEmail()))
	if email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if len(email) > 254 {
		return status.Error(codes.InvalidArgument, "email is too long")
	}
	if !emailRegex.MatchString(email) {
		return status.Error(codes.InvalidArgument, "email is invalid")
	}

	// --- PASSWORD ---
	password := req.GetPassword()
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if len(password) < 8 {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}
	if len(password) > 72 {
		return status.Error(codes.InvalidArgument, "password is too long")
	}

	if strings.Contains(strings.ToLower(password), email) {
		return status.Error(codes.InvalidArgument, "password must not contain email")
	}

	if !isStrongPassword(password) {
		return status.Error(
			codes.InvalidArgument,
			"password must contain uppercase, lowercase, number and special character",
		)
	}

	return nil
}

func isStrongPassword(pw string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, c := range pw {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~", c):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
