package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestMakeJWT(t *testing.T) {
	uuid := uuid.MustParse("e3aff3f2-fd5b-4390-b4e8-e5c861266cb6")
	tokenSecret := "really secret!!!"
	expiresIn := time.Minute

	signedToken, err := MakeJWT(uuid, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("error while signing token")
	} else {
		fmt.Println(signedToken)
	}
}


