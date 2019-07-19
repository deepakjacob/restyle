package signer

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	jwtKey          = "JWT"
	upProject       = "Up Project"
	jwtSharedKey    = "JWT_SHARED_KEY"
	jwtSharedEncKey = "JWT_SHARED_ENC_KEY"
)

var (
	sharedKey = []byte(os.Getenv(jwtSharedKey))
	// JWT_SHARED_ENC_KEY should be 16 bytes long. ie, 16 characters
	sharedEncryptionKey = []byte(os.Getenv(jwtSharedEncKey))
	signer, _           = jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256, Key: sharedKey},
		&jose.SignerOptions{})
	rsaSigner = mustMakeSigner(jose.RS256, rsaPrivKey)
	enc, _    = jose.NewEncrypter(
		jose.A128GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       sharedEncryptionKey,
		},
		(&jose.EncrypterOptions{}).WithType(jwtKey).WithContentType(jwtKey),
	)

	errEmptyVals = errors.New("input email or id is empty")
)

//SignEncryptJWT sign and encrypt JWTs with claims
func SignEncryptJWT(user *domain.UserToken) (string, error) {
	if user.Email == "" || user.UserID == "" {
		logger.Log.Error("jwt signing failed", zap.Error(errEmptyVals))
		return "", errEmptyVals
	}

	cl := jwt.Claims{
		Subject: user.Email,
		ID:      user.UserID,
		Issuer:  upProject,
	}

	raw, err := jwt.
		SignedAndEncrypted(rsaSigner, enc).
		Claims(cl).
		CompactSerialize()

	if err != nil {
		logger.Log.Error("jwt signing failed",
			zap.String("UserID", user.UserID),
			zap.Error(err))
		return "", err
	}

	return raw, nil
}

// DecryptJWT decrypt JWT represented as a raw string
func DecryptJWT(raw string) (*domain.UserToken, error) {
	tok, err := jwt.ParseSignedAndEncrypted(raw)
	if err != nil {
		logger.Log.Error("jwt parsing failed", zap.Error(err))
		return nil, err
	}

	nested, err := tok.Decrypt(sharedEncryptionKey)
	if err != nil {
		logger.Log.Error("jwt decryption failed", zap.Error(err))
		return nil, err
	}

	claims := jwt.Claims{}
	if err := nested.Claims(&rsaPrivKey.PublicKey, &claims); err != nil {
		logger.Log.Error("jwt claims extraction failed", zap.Error(err))
		return nil, err
	}

	return &domain.UserToken{Email: claims.Subject}, nil
}

func mustMakeSigner(alg jose.SignatureAlgorithm, k interface{}) jose.Signer {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: alg, Key: k}, nil)
	if err != nil {
		panic("failed to create signer:" + err.Error())
	}

	return sig
}

func mustUnmarshalRSA(data string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		panic("failed to decode PEM data")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse RSA key: " + err.Error())
	}
	if key, ok := key.(*rsa.PrivateKey); ok {
		return key
	}
	panic("key is not of type *rsa.PrivateKey")
}

var rsaPrivKey = mustUnmarshalRSA(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDUacqfGf6oNtzl
HvHLo7Im8pbGKPvvmavjvQEFLwM55UiYurafT5FBo1VWwBqCkvZgJ3HlPFQQYlWe
KPE9Utcxtl99+xD4yEJHf1qKGD+J+fs2E5+h/lJmRkiElCkg+SGZ+MTald5RO0iO
Inru2NM/XhHxFjvFjpWapoj4Yn3oMiXx9lK/LJZzvWzE/V2ZKxQwo7rT9C/cW4St
vTbj39JjtmOc1X/biLsut8jHNGL8PU3iMqEKvgmrS5SC5I84EZlO1PKFAO4ad4w+
wQgPqkqgwTagE1drpe+ZNK/2lpQ835DXK/XBouffwxVopAaxZj0NGQXjFhPK3NhD
hQN3NikhAgMBAAECggEAfT13E4HvLV1TtuAc7tWwp9gm3+WwTeBMgfWhkV3byBoR
SMPmA16VpZ7ZJVIVD2H5VE5NkDyW0CY6lwtK2b5rVUtTWGNc6WKh+af0STHs3LyA
yqydVZUvHlBYV5tH+MILds4/uyXcNAcBS5S26Pb87wLXD0/tpBjiypFgdxsTY4CB
+1wjzjabPkxKEFEHtQXL0o3+0j/lZA16GA2pt3BU5Oilqog72D+hwT94aE2E/F4O
jjW1JZ3Wxj3V5urSy3aPG+P60VlmUn5ir8KIRbDMrEuZPXvGjp2hbj3NjkikTetX
VnWax1QnLtdrybK8xz3fG32ak8VG7NNbS2zYDml3MQKBgQD9qiNUJNr8GaU9QjU9
J/acqel+GJKKNhTmKmxantoiItH29116RvvRVOYSyx5xcZNuUZ2RxPVgNKEwLigr
cBr2pBHoaR9oaWsTwdA4RmZvlfa2AX3tzeiD2y7uH+vU+gHBJsykVKQbBCUVxKWv
vVuDpZJGme6P71pLqGMDemsdPQKBgQDWXm2dyb66dS5RTtw9hqdVbLSR/sYQ3snq
CRJEo/xkK10hMtyzZJM1EPcE3wI5sIn19K+C4oUu0BEt39nu1UpurwIKFIn169kq
L2fkCaCXWKx7CLKTcBC5cuLFfQc6vGEol6J355szr5bn1ElzpYzYYzVuMTSkMZ74
VoXB4CZBtQKBgE69HEBHPG5aq48LWSlFmWhh2aeZiws55FzZuE6c1osYEeK+QBbv
p8T/vHcC8800+xWYYffYkm4tiAdDnJZ3MvdyUFi7INOxPVqho7eEKtHiU/WFGwjR
DKa5R6UE5Zhzjk3ddFJiL5pvO++43dFiTuDbaT9fEs089+NaPnna6xrlAoGAR0Vb
0nMJu3pMLPI4HSiQp8Edg9Cdz2wS24GqljGjLzaPnwMHB4mvu5vpVLBEUCPWqnRw
ieZ/+yFoJMVg8pvtREFhPzK275E7QWBDfTiKMOdlaP5qSMSgetesd5Zq+ec5skI/
3PeezR9a37bfuNhYrHTHhdxMMS7iOJSjoeLBNWUCgYEAueDlwnxfu7jwO+BrDXmr
d4dYM+UjMXFu0uj60r/qll2FMtDdCvzQN5wAIO8D53uAklb1E85wk3SenNtfFTzf
po5YR/jL0l1QvPhyf52YFwExEpB7TmTO413XzuLGSfWu1MlgabfoXRxcqJaIOo/V
MD+UC9vAtSd/CogO5fF+i0E=
-----END PRIVATE KEY-----`)
