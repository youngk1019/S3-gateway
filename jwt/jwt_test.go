package jwt

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestJWT(t *testing.T) {
	claims, err := ParseJWT("eyJhbGciOiJIUzUxMiJ9.eyJrZXlfdV9uYW1lIjoiNzkyNjAyNTIwQHFxLmNvbSIsImtleV91aWQiOiJmMzg0MmI2OTQ4YzM0MDYzODU2ODNkODIzNmI4MzExMiIsImxvZ2luX3VzZXJfa2V5IjoiMjQ4YmQzYmEtZjFkNC00MWNiLWIzNjUtZGQyMDNlNjMyNjkzIn0.j4tRe000k8zDjpLN7UY9kACO-OS-1cKVIh-Egy-j714CIlJp8Utfq2c-mX4-Nsr8h0RI_VuMqgKLoO0bHJKIkg")
	assert.Equal(t, err, nil)
	assert.Equal(t, claims["key_u_name"], "792602520@qq.com")
	assert.Equal(t, claims["key_uid"], "f3842b6948c3406385683d8236b83112")
	assert.Equal(t, claims["login_user_key"], "248bd3ba-f1d4-41cb-b365-dd203e632693")
	assert.Equal(t, claims["test"], nil)

	claims, err = ParseJWT("")
	assert.Equal(t, err == nil, false)
	assert.Equal(t, claims["key_u_name"], nil)
	assert.Equal(t, claims["key_uid"], nil)
	assert.Equal(t, claims["login_user_key"], nil)
	assert.Equal(t, claims["test"], nil)
}
