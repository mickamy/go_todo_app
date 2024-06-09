package auth

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/mickamy/go_todo_app/clock"
	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/store"
	"github.com/mickamy/go_todo_app/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("embedded public key does not contain expected bytes want=(%s) got=(%s)", want, rawPubKey)
	}
	want = []byte("-----BEGIN PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("embedded private key does not contain expected bytes want=(%s) got=(%s)", want, rawPrivKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	ctx := context.Background()
	moq := &StoreMock{}
	wantID := entity.UserID(20)
	u := fixture.User(&entity.User{ID: wantID})
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("user id want=%d got=%d", wantID, userID)
		}
		return nil
	}
	sut, err := NewJWTer(moq, clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := sut.GenerateToken(ctx, *u)
	if err != nil {
		t.Fatalf("GenerateToken got err=%v, want err=nil", err)
	}
	if len(got) == 0 {
		t.Errorf("GenerateToken got len(got)=0")
	}
}

func TestJWTer_GetToken(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}
	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("github.com/mickamy/go_todo_app").
		Subject("access_token").
		IssuedAt(c.Now()).
		// Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test_role").
		Claim(UserNameKey, "test_user_name").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	pkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}
	userID := entity.UserID(20)

	ctx := context.Background()
	moq := &StoreMock{}
	moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
		return userID, nil
	}
	sut, err := NewJWTer(moq, c)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "https://github.com/mickamy", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(signed)))
	got, err := sut.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("GetToken got err=%v, want err=nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetToken got=%v, want %v", got, want)
	}
}

type FixedTomorrowClocker struct{}

func (c FixedTomorrowClocker) Now() time.Time {
	return clock.FixedClocker{}.Now().Add(24 * time.Hour)
}

func TestJWTer_GetToken_NG(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}
	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("github.com/mickamy/go_todo_app").
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test_role").
		Claim(UserNameKey, "test_user_name").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	pkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}

	type moq struct {
		userID entity.UserID
		err    error
	}
	tests := map[string]struct {
		c   clock.Clocker
		moq moq
	}{
		"expire": {
			c: FixedTomorrowClocker{},
		},
		"notFoundInStore": {
			c: clock.FixedClocker{},
			moq: moq{
				err: store.ErrNotFound,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			moq := &StoreMock{}
			moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
				return tt.moq.userID, nil
			}
			sut, err := NewJWTer(moq, tt.c)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodGet, "https://github.com/mickamy", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(signed)))
			got, err := sut.GetToken(ctx, req)
			if err == nil {
				t.Errorf("GetToken got nil, want error")
			}
			if got != nil {
				t.Errorf("GetToken got %v, want nil", got)
			}
		})
	}
}
