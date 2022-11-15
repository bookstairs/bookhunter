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
)

// login is used first time you execute the command line.
func (t *Telegram) login() error {
	return t.execute(func(_ context.Context, _ *telegram.Client) error {
		return nil
	})
}

// authentication is used for log into the telegram with a session support.
// Every telegram execution will require this method.
func (t *Telegram) authentication(ctx context.Context) error {
	// Setting up authentication flow helper based on terminal auth.
	flow := auth.NewFlow(&terminalAuth{mobile: t.mobile}, auth.SendCodeOptions{})
	if err := t.client.Auth().IfNecessary(ctx, flow); err != nil {
		return err
	}

	status, _ := t.client.Auth().Status(ctx)
	if !status.Authorized {
		return errors.New("failed to login, please check you login info or refresh the session by --refresh")
	}

	return nil
}

// terminalAuth implements authentication via terminal.
type terminalAuth struct {
	mobile string
}

func (t *terminalAuth) Phone(_ context.Context) (string, error) {
	// Make the mobile number has the country code as the prefix.
	addCountryCode := func(mobile string) string {
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

	if t.mobile == "" {
		fmt.Print("Enter Phone Number (+86): ")
		phone, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}
		t.mobile = addCountryCode(phone)
	}

	return t.mobile, nil
}

func (t *terminalAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := term.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (t *terminalAuth) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

func (t *terminalAuth) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("signup call is not expected")
}

func (t *terminalAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}
