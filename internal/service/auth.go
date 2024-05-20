package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/Futturi/Raspisanie/internal/repository"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	salt1 = "rwpkjgop0pqwpitgnvc,xz;aqwperngbxlsprfvmdlsleprfogjbm,cla'apfkvnfjksmvmrkdc"
	salt2 = "woegjmpoeqpwetuyturgnvmc,x,za;dlfkghgitorpepdkcmv bngmfls;zpaiotjgmb"
)

type AuthService struct {
	repo repository.Auth
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) SignUp(group entities.User) (int, error) {
	gr := entities.User{
		Email:    group.Email,
		Password: hashPass(group.Password),
		Group:    group.Group,
		Name:     group.Name,
	}
	return a.repo.SignUp(gr)
}

func hashPass(pass string) string {
	hash := sha1.New()
	hash.Write([]byte(pass))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt1)))
}

type Claims struct {
	Id    int    `json:"id"`
	Group string `json:"group"`
	jwt.StandardClaims
}

func (a *AuthService) SignIn(user entities.User) (string, error) {
	us := entities.User{
		Email:    user.Email,
		Password: hashPass(user.Password),
	}
	id, group, err := a.repo.SignIn(us)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{id, group, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}})

	return token.SignedString([]byte(salt2))
}

func (a *AuthService) ParseToken(header string) (int, string, error) {
	token, err := jwt.ParseWithClaims(header, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, errors.New("invalid signing method")
		}
		return []byte(salt2), nil
	})
	if err != nil {
		return 0, "", err
	}
	Claims1, ok := token.Claims.(*Claims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}
	return Claims1.Id, Claims1.Group, nil
}
