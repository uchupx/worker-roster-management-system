package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessToken struct {
	Data interface{} `json:"data"`
	Exp  int64       `json:"exp"`
}

type RefreshToken struct {
	Id  string `json:"id"`
	Exp int64  `json:"exp"`
}

type JWTConf struct {
	Privatekey string
	PublicKey string
}

type CryptService interface {
	CreateSignPSS(value string) (signatureStr string, err error)
	Verify(value string, signature string) (resp bool, err error)
	CreateJWTToken(expDuration time.Duration, content interface{}) (token *string, err error)
	CreateRefreshToken(expDuration time.Duration, id string) (token *string, err error)
	CreateAccessToken(expDuration time.Duration, content interface{}) (token *string, err error)
	VerifyJWTToken(token string) (result interface{}, err error)
}

type cryptService struct {
	conf JWTConf
}

type Params struct {
	Conf JWTConf
}

func NewCryptService(params Params) CryptService {
	return &cryptService{
		conf: params.Conf,
	}
}

func (s *cryptService) loadRsaPrivateKey() (rsaKey *rsa.PrivateKey, err error) {
	key := strings.ReplaceAll(s.conf.Privatekey, "\\n", "\n") // remove double slash, because it will affected when convert to byte
	rsaKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		// s.logger.Errorf("[loadRsaPrivateKey] failed parse private key, err: %+v", err)
		return
	}

	return
}

func (s *cryptService) loadRsaPublicKey() (rsaPub *rsa.PublicKey, err error) {
	key := strings.ReplaceAll(s.conf.PublicKey, "\\n", "\n") // remove double slash, because it will affected when convert to byte
	rsaPub, err = jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		// s.logger.Errorf("[loadRsaPublicKey] failed parse public key, err: %+v", err)
		return
	}

	return
}

func (s *cryptService) CreateSignPSS(value string) (signatureStr string, err error) {
	msg := []byte(value)
	msgHash := sha256.New()

	_, err = msgHash.Write(msg)
	if err != nil {
		return
	}

	msgHashSum := msgHash.Sum(nil)

	privateKey, err := s.loadRsaPrivateKey()
	if err != nil {
		return
	}

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		// s.logger.Errorf("[CreateSignPSS] failed to create signature, err: %+v", err)
		return
	}

	signatureStr = base64.URLEncoding.EncodeToString(signature)
	return
}

func (s *cryptService) Verify(value string, signature string) (resp bool, err error) {
	resp = false

	signatureByte, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return
	}

	msg := []byte(value)
	msgHash := sha256.New()

	_, err = msgHash.Write(msg)
	if err != nil {
		return
	}

	msgHashSum := msgHash.Sum(nil)

	publicKey, err := s.loadRsaPublicKey()
	if err != nil {
		return
	}

	err = rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSum, signatureByte, nil)
	if err != nil {
		return
	}

	resp = true

	return
}

func (s *cryptService) CreateJWTToken(expDuration time.Duration, content interface{}) (token *string, err error) {
	privateKey, err := s.loadRsaPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("[CreateJWTToken] failed load rsa private key, err: %+v", err)
	}

	publicKey, err := s.loadRsaPublicKey()
	if err != nil {
		return nil, fmt.Errorf("[CreateJWTToken] failed load rsa public key, err: %+v", err)
	}

	jwtServicecryptService := NewJWT(privateKey, publicKey)

	return jwtServicecryptService.Create(expDuration, content)
}

func (s *cryptService) CreateRefreshToken(expDuration time.Duration, id string) (token *string, err error) {
	exp := time.Now().Add(expDuration).Unix()

	refreshToken := &RefreshToken{
		Id:  id,
		Exp: exp,
	}

	return s.CreateJWTToken(expDuration, refreshToken)
}

func (s *cryptService) CreateAccessToken(expDuration time.Duration, content interface{}) (token *string, err error) {
	accessToken := &AccessToken{
		Data: content,
		Exp:  time.Now().Add(expDuration).Unix(),
	}

	return s.CreateJWTToken(expDuration, accessToken)
}

func (s *cryptService) VerifyJWTToken(token string) (result interface{}, err error) {
	privateKey, err := s.loadRsaPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("[VerifyJWTToken] failed load rsa private key, err: %+v", err)
	}

	publicKey, err := s.loadRsaPublicKey()
	if err != nil {
		return nil, fmt.Errorf("[VerifyJWTToken] failed load rsa public key, err: %+v", err)
	}

	jwtServicecryptService := NewJWT(privateKey, publicKey)

	return jwtServicecryptService.Validate(token)
}
