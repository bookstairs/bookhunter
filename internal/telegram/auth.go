package telegram

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"golang.org/x/term"

	"github.com/bookstairs/bookhunter/internal/log"
)

// TerminalAuth implements authentication via terminal.
type TerminalAuth struct {
	mobile string
}

func NewAuth(mobile string) auth.UserAuthenticator {
	return &TerminalAuth{mobile: mobile}
}

func (t *TerminalAuth) Phone(_ context.Context) (string, error) {
	if t.mobile == "" {
		fmt.Print("Enter Phone Number (+86): ")
		phone, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}
		t.mobile = AddCountryCode(phone)
	}
	return t.mobile, nil
}

func (t *TerminalAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := term.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (t *TerminalAuth) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

func (t *TerminalAuth) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("signup call is not expected")
}

func (t *TerminalAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

// AddCountryCode will return the mobile number with country code as the prefix.
func AddCountryCode(mobile string) string {
	if strings.HasPrefix(mobile, "+") {
		return mobile
	} else if strings.HasPrefix(mobile, "86") {
		return "+" + mobile
	} else if mobile != "" {
		return "+86" + mobile
	} else {
		return ""
	}
}

// login is used for log into the telegram with a session support.
func login(ctx context.Context, client *telegram.Client, config *Config) error {
	// Setting up authentication flow helper based on terminal auth.
	flow := auth.NewFlow(NewAuth(config.Mobile), auth.SendCodeOptions{})
	if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		log.Fatal(err)
	}

	status, _ := client.Auth().Status(ctx)
	if !status.Authorized {
		return fmt.Errorf("failed to login, please check you login info or refresh the session by --refresh")
	}

	return nil
}
