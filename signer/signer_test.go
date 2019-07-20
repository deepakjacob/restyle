package signer

import (
	"log"
	"testing"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
)

func init() {
	if err := logger.Init(-1, "2006-01-02T15:04:05Z07:00"); err != nil {
		log.Fatal(err)
	}
}
func TestEncryptDecrypt(t *testing.T) {
	user := &domain.UserToken{
		UserID: "1234577979",
		Email:  "test@test.com",
	}
	signer := &Signer{}

	got, err := signer.SignEncrypt(user)
	if err != nil {
		t.Error("Error while signing/ encrypting value")
	}
	// fmt.Println("Raw JWT string =>", got)

	want, err := signer.Decrypt(got)
	if err != nil {
		t.Error("error while decrypting value")
	}
	if user.Email != want.Email {
		t.Errorf("unexpected value \n got: %s \n want: %v", got, want)
	}
}

var (
	result string
)

func TestDecrypt(t *testing.T) {
	raw := `eyJhbGciOiJkaXIiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTI4R0NNIiwidHlwIjoiSldUIn0..d4172A3ckbFBC_QB.i_uv2vTZ7lnpyba4XxELUiPNm1iDSGbZlN6zeRsojxqbeXAZToie2f_g2x8Vqjal9RVSEvzHqoJyOWHIupX4e3JKVXByy1bCGaxp_LqJD7VZ04itkgD5CzDf-9trVxRglKVUWU7H2OKkAICeRtUth57JBZbqZ2OrkyHMKTRS73FZeSQndhDUc6GhdGiHlBtMADsuZPduEx4c3mgTvnM2W2f9VuEvAr9CiuC4g-Q9M06tytNuJYBZPpYO7gVwV8wGGLVH-K-0tE7LQwiKVjXTgHdBHH5zSxHXU34gOnqf93phnT2LrNNzE0vfadMt2D1tgXeYe2Mu0FvmRX-9jyEG5p_EFBJmb6owA3wSDdOXsuSUrzPyQ9ZONTPBz4wnOxfd1fqNY25jSWmOKaSwR6IlXwib9OkCJHXzDQ5pZkzx3SSo0NVe7fcv7oF9O2spoOlrcoKu92aO5OZ_jGq6EfbwndFHVLqyaPe-IcNbPOg960yibk4K0II0VDOUK0whopZva_d1WnR_rkc9C-N8r0UWL_qSz1qDDV538-xDljQWwXgG_dljlsLwI7vaADN7Q53m.3NV5BH1CoNvLMsBkZ_EfXQ`
	signer := &Signer{}
	got, err := signer.Decrypt(raw)
	if got.Email != "test@test.com" {
		t.Errorf("Error while decrypting value %+v", err)
	}

}

func BenchmarkEncrypt(b *testing.B) {
	user := &domain.UserToken{
		UserID: "1234577979",
		Email:  "test@test.com",
	}
	var got string
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		signer := &Signer{}
		got, _ = signer.SignEncrypt(user)
	}
	result = got
}
func BenchmarkDecrypt(b *testing.B) {
	raw := `eyJhbGciOiJkaXIiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTI4R0NNIiwidHlwIjoiSldUIn0..d4172A3ckbFBC_QB.i_uv2vTZ7lnpyba4XxELUiPNm1iDSGbZlN6zeRsojxqbeXAZToie2f_g2x8Vqjal9RVSEvzHqoJyOWHIupX4e3JKVXByy1bCGaxp_LqJD7VZ04itkgD5CzDf-9trVxRglKVUWU7H2OKkAICeRtUth57JBZbqZ2OrkyHMKTRS73FZeSQndhDUc6GhdGiHlBtMADsuZPduEx4c3mgTvnM2W2f9VuEvAr9CiuC4g-Q9M06tytNuJYBZPpYO7gVwV8wGGLVH-K-0tE7LQwiKVjXTgHdBHH5zSxHXU34gOnqf93phnT2LrNNzE0vfadMt2D1tgXeYe2Mu0FvmRX-9jyEG5p_EFBJmb6owA3wSDdOXsuSUrzPyQ9ZONTPBz4wnOxfd1fqNY25jSWmOKaSwR6IlXwib9OkCJHXzDQ5pZkzx3SSo0NVe7fcv7oF9O2spoOlrcoKu92aO5OZ_jGq6EfbwndFHVLqyaPe-IcNbPOg960yibk4K0II0VDOUK0whopZva_d1WnR_rkc9C-N8r0UWL_qSz1qDDV538-xDljQWwXgG_dljlsLwI7vaADN7Q53m.3NV5BH1CoNvLMsBkZ_EfXQ`
	var got *domain.UserToken
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		signer := &Signer{}
		got, _ = signer.Decrypt(raw)
	}
	result = got.UserID
}
