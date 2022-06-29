package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zibbp/ganymede/ent"
	entUser "github.com/zibbp/ganymede/ent/user"
	"github.com/zibbp/ganymede/internal/database"
	"github.com/zibbp/ganymede/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Store *database.Database
}

func NewService(store *database.Database) *Service {
	return &Service{Store: store}
}

func (s *Service) Register(c echo.Context, user user.User) (*ent.User, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	u, err := s.Store.Client.User.Create().SetUsername(user.Username).SetPassword(string(hashedPassword)).Save(c.Request().Context())
	if err != nil {
		if _, ok := err.(*ent.ConstraintError); ok {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return u, nil
}

func (s *Service) Login(c echo.Context, uDto user.User) (*ent.User, error) {
	u, err := s.Store.Client.User.Query().Where(entUser.Username(uDto.Username)).Only(c.Request().Context())
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(uDto.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	uDto = user.User{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}

	// Generate JWT and set cookie
	err = GenerateTokensAndSetCookies(&uDto, c)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %v", err)
	}

	return u, nil
}

func (s *Service) Refresh(c echo.Context, refreshToken string) error {

	tkn, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetJWTRefreshSecret()), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return fmt.Errorf("invalid refresh token")
		}
		return fmt.Errorf("error parsing refresh token: %v", err)
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		uID := claims["user_id"].(string)
		uUUID, err := uuid.Parse(uID)
		if err != nil {
			return fmt.Errorf("error parsing user id: %v", err)
		}
		u, err := s.Store.Client.User.Query().Where(entUser.ID(uUUID)).Only(c.Request().Context())
		if err != nil {
			return fmt.Errorf("error getting user: %v", err)
		}

		// Generate JWT and set cookie
		err = GenerateTokensAndSetCookies(&user.User{ID: u.ID, Username: u.Username, Role: u.Role}, c)
		if err != nil {
			return fmt.Errorf("error generating tokens: %v", err)
		}
		return nil
	}

	return err
}