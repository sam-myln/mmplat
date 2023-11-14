package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"mmplat/internal/util"
	"strings"
)

func newOtpCmd(ctx *CmdCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "otp",
		Short: "Generates otp image",
		Long:  "Cool stuff",
		Example: `otp --user anton:anton --login-data FILE OR
otp --user anton:anton --login-data FILE...`,
		RunE: ctx.OtpCmdRunE,
	}
	cmd.PersistentFlags().StringP("user", "u", "", "user data in format: user:pass")
	cmd.MarkFlagRequired("user")

	return cmd
}

func (ctx *CmdCtx) OtpCmdRunE(self *cobra.Command, _ []string) error {
	loginPass, _ := ctx.flags.GetStringSlice(cmdFlagNameLoginData)
	credMngr := util.CreateCredMngr(loginPass...)
	credMngr.ParseCredentials()
	t, _ := self.Flags().GetString("user")
	t1 := strings.Split(t, ":")
	if !credMngr.Exists(t1[0],t1[1]) {
		ctx.log.Errorf("error: otp: user not found, or parsing error")
		return errors.New("user not found")
	}
	key, err := util.KeyGen(t)
	if err != nil {
		ctx.log.Errorf("error: otp: keygen failed w/: %v", err)
		return err
	}
	fmt.Printf("User: %s\nSecret: %s\n", t1[0], key.Secret())
	return nil
}
