/*-----------------------------------------------------------------------------
# Name:        NTLM_AUTH Secure Wrapper 0.0.4
# Purpose:     Interact with NTLM_AUTH (samba) to provide safer authorization
#			  	for apache mod_authnz_external (or similar)
#
# Author:      Rafal Wilk <rw@pcboot.pl>
#
# Created:     08-06-2022
# Modified:    09-06-2022
# Copyright:   (c) PcBoot 2022
# License:     BSD-new
-----------------------------------------------------------------------------*/

package main

import (
	"bufio"
	"fmt"
	"github.com/alexflint/go-arg"
	"os"
	"os/exec"
	"strings"
)

const VERSION = "0.0.4"

var args struct {
	NTLMAuth   string `arg:"-e,--ntlm-auth-binary,required" help:"NTLM_AUTH binary"`
	Domain     string `arg:"-d,--domain" help:"domain name"`
	Membership string `arg:"-m,--require-membership-of" help:"require that a user be a member of this group"`
	VersionPrn bool   `arg:"-v,--version" help:"print version"`
}

func main() {
	if err := arg.Parse(&args); err != nil {
		printHeader()
		arg.MustParse(&args)
	}

	if args.VersionPrn {
		printHeader()
		return
	}

	reader := bufio.NewReader(os.Stdin)
	user, _ := reader.ReadString('\n')
	pass, _ := reader.ReadString('\n')

	user = strings.Trim(user, "\n")
	pass = strings.Trim(pass, "\n")

	isOk, err := verifyCredentials(user, pass)
	if err != nil {
		panic(err)
	}
	if isOk {
		fmt.Println("OK")
	} else {
		fmt.Println("WRONG_CREDENTIALS")
		os.Exit(1)
	}

}

func verifyCredentials(username string, password string) (bool, error) {
	var cargs []string
	cargs = append(cargs, fmt.Sprintf("--username=%s", username))

	if len(args.Domain) > 0 {
		cargs = append(cargs, fmt.Sprintf("--domain=%s", args.Domain))
	}

	if len(args.Membership) > 0 {
		cargs = append(cargs, fmt.Sprintf("--require-membership-of=%s", args.Membership))
	}

	cmd := exec.Command(args.NTLMAuth, cargs...)

	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s\n", password))
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		if err.Error() != "exit status 1" {
			return false, err
		}
	}

	if strings.Contains(string(stdoutStderr), "NT_STATUS_OK") {
		return true, nil
	}

	return false, nil
}

func printHeader() {
	fmt.Printf("NTLM_AUTH Secure Wrapper %s by <rw@pcboot.pl>\n", VERSION)
	fmt.Println("All rights reserved. (c) PcBoot 2022")
	fmt.Println()
}
