package pixiv

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/dghubble/sling"
	"golang.org/x/xerrors"
)

// default oauth client params
const (
	DefaultClientID     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	DefaultClientSecret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	DefaultClientHash   = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
)

type authParams struct {
	GetSecureURL int    `url:"get_secure_url,omitempty"`
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	Username     string `url:"username,omitempty"`
	Password     string `url:"password,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
}

type loginResponse struct {
	Response *authInfo `json:"response"`
}

type authInfo struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Scope        string   `json:"scope"`
	RefreshToken string   `json:"refresh_token"`
	User         *Account `json:"user"`
	DeviceToken  string   `json:"device_token"`
}

type loginError struct {
	HasError bool                 `json:"has_error"`
	Errors   map[string]authError `json:"errors"`
}

type authError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Tokens from auth login
type Tokens struct {
	AccessToken   string
	RefreshToken  string
	tokenDeadline time.Time
}

// Session for pixiv
type Session struct {
	ClientID     string
	ClientSecret string
	ClientHash   string

	Account *Account

	tokens *Tokens
	r      *sling.Sling
}

// NewSession for pixiv
func NewSession() *Session {
	s := new(Session)
	s.ClientID = DefaultClientID
	s.ClientSecret = DefaultClientSecret
	s.ClientHash = DefaultClientHash

	s.r = sling.New().
		Base("https://app-api.pixiv.net/").
		Set("User-Agent", "PixivAndroidApp/5.0.64 (Android 6.0)").
		Set("App-Version", "6.7.1").
		Set("App-OS-VERSION", "10.3.1").
		Set("App-OS", "ios")

	return s
}

// Login into pixiv for private api
func (s *Session) Login(username, password string) (*Tokens, error) {
	if err := s.auth(&authParams{
		Username: username,
		Password: password,
	}); err != nil {
		return nil, err
	}
	return s.tokens, nil
}

// Tokens returns the different tokens
func (s *Session) Tokens() *Tokens {
	return s.tokens
}

// RefreshLogin with the given refresh_token
func (s *Session) RefreshLogin(tokens *Tokens) error {
	s.tokens = tokens
	return s.refreshAuth()
}

func (s *Session) auth(params *authParams) error {
	ts := time.Now().Format(time.RFC3339)
	hasher := md5.New()
	hasher.Write([]byte(ts))
	hasher.Write([]byte(s.ClientHash))
	hash := hex.EncodeToString(hasher.Sum(nil))

	slinger := sling.New().
		Base("https://oauth.secure.pixiv.net/").
		Set("User-Agent", "PixivAndroidApp/5.0.64 (Android 6.0)").
		Set("X-Client-Time", ts).
		Set("X-Client-Hash", hash)

	params.GetSecureURL = 1
	params.ClientID = s.ClientID
	params.ClientSecret = s.ClientSecret

	if params.GrantType == "" {
		if params.RefreshToken != "" {
			params.GrantType = "refresh_token"
		} else {
			params.GrantType = "password"
		}
	}

	res := &loginResponse{
		Response: &authInfo{
			User: &Account{},
		},
	}
	loginErr := &loginError{
		Errors: map[string]authError{},
	}

	_, err := slinger.New().Post("auth/token").BodyForm(params).Receive(res, loginErr)
	if err != nil {
		return err
	}
	if loginErr.HasError {
		for k, v := range loginErr.Errors {
			return xerrors.Errorf("error [%v] %v : %w", k, v.Message, ErrAuthentication)
		}
	}

	s.tokens = &Tokens{
		AccessToken:   res.Response.AccessToken,
		RefreshToken:  res.Response.RefreshToken,
		tokenDeadline: time.Now().Add(time.Duration(res.Response.ExpiresIn) + time.Second),
	}

	return nil
}

func (s *Session) refreshAuth() error {
	s.r = s.r.Set("Authorization", "Bearer "+s.tokens.AccessToken)
	if time.Now().Before(s.tokens.tokenDeadline) {
		return nil
	}

	if s.tokens.RefreshToken == "" {
		return ErrMissingToken
	}

	params := &authParams{
		RefreshToken: s.tokens.RefreshToken,
	}

	return s.auth(params)
}
