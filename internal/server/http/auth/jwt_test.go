package auth

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todoNote/internal/model"
)

const (
	private="LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ4d2dnRTdBZ0VBQWtFQXZTS1B0WnFwdngxQ0lsNDgKRjdrd3dEN2w4dCs5M3VMMXc0OVNZOGhXQjJLdEFFemg3d2lodU1UZ3drYW1ESHRZdnpNaTQrK2w0bGdNTXBCUQpZMDVxNFFJREFRQUJBa0JFVDFWRUxBWWU1bnhhV1ZxdTNzNEN3VFRnRVh0TUl3RE1qdGtjL09CRmJmZTFISERJCkllZ1VnQWdFSk5xclBZVDhvcEZPMGI4UmpqKytzamNxS05BaEFpRUEvVHJhck82ZXNNdGhSMENsRkVmaElXSEIKSEVHbDR1Y0tzaTRob1Q2VnhUMENJUUMvTkRjekVmNEFXSHRTcjNVUTd5WG52WlBLeWN2T2xWSU5sZXZYV2RpKwpkUUloQUxEN3JWSW9CQ2swTyt6OHRXT1RTVGwzaE93bXhiWHNISUdqMUVWSjVJdFJBaUEzZWVmMkttYy9KRzBMCnJaclN3Z0NHZjR2TkQ4WFJkNk9xQzNDMU4vMWFMUUloQU9SanVMM1N1T3lJZjNpNi95RDg0YVdMSzdGcGo2MjYKOElabVRLSU9FZWNPCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0="
	public="LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTDBpajdXYXFiOGRRaUplUEJlNU1NQSs1ZkxmdmQ3aQo5Y09QVW1QSVZnZGlyUUJNNGU4SW9iakU0TUpHcGd4N1dMOHpJdVB2cGVKWURES1FVR05PYXVFQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
)

func TestCreateValidateToken(t *testing.T) {
	tts := []struct {
		id      model.Id
		expires time.Duration
		wait   time.Duration
		hasErr bool
	}{
		{1, time.Second * 2, time.Second * 3, true},
		{2, time.Second * 3, time.Second * 2, false},
	}
	for _, tt := range tts {
		t.Run(fmt.Sprintf("test %v", tt.id), func(t *testing.T) {
			auth, err := NewJwtAuth(tt.expires, private, public)
			if err != nil {
				t.Fatal(err)
			}

			token, err := auth.CreateToken(model.UserInReq{Id: tt.id})
			time.Sleep(tt.wait)

			u, err := auth.ValidateToken(token)

			actual := err != nil
			assert.Equal(t, tt.hasErr, actual)
			if !tt.hasErr {
				assert.EqualValues(t, tt.id, u.Id)
			}
		})
	}
}
