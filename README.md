
# NTLM_AUTH Secure Wrapper

NTLM_AUTH Secure Wrapper interacts with NTLM_AUTH (samba) to provide safer authorization for apache mod_authnz_external (or similar).

## Background

If you're trying to set up Apache 2.4 Basic authorization based on NTLM_AUTH probably you found major security risk 
around executing standard ntlm_auth with username and password as arguments. For production use it's unacceptable.
NTLM_AUTH Secure Wrapper provides simple solution - compatible with "pipe method" interface for ntlm_auth. **Wrapper gets
username and password as separate input and executes ntlm_auth without passing password as argument.**

### Why does it matter?

...because executing process with password as argument is huge risk - the password will be visible in `ps -e` output. 
NTLM_AUTH Secure Wrapper solves this problem.

## Getting started

A quick introduction of the minimal setup you need to get started.

```shell
wget https://github.com/rwilk/ntlm_auth_wrapper/raw/master/bin/linux/ntlm_auth_wrapper # linux
[or]
wget https://github.com/rwilk/ntlm_auth_wrapper/raw/master/bin/freebsd/ntlm_auth_wrapper # freebsd

./ntlm_auth_wrapper -e /usr/bin/ntlm_auth -d MYDOMAIN.local 
```

If you want to limit authorized users to specified AD Security Group, you can also use *-m* argument:

```shell
./ntlm_auth_wrapper -e /usr/bin/ntlm_auth -d MYDOMAIN.local -m MYDOMAIN\\Required-AD-Group
```

### Compiling

If you need version for other platform or just want to compile code yourself you can use below commands:

```shell
git clone https://github.com/rwilk/ntlm_auth_wrapper.git
cd ntlm_auth_wrapper/
go build
```

### Example usage case

Below example apache24 virtual host config.

**Make sure that `ntlm_auth` and `mod_authnz_external` are installed, configured and enabled**, then:

```apacheconf

<VirtualHost *:443>
    ServerAdmin admin@example.com
    ServerName example.com
    
    [...]

    DefineExternalAuth ntlmAuth pipe "/tools/ntlm_auth_wrapper/ntlm_auth_wrapper -e /usr/local/bin/ntlm_auth -m MYDOMAIN.local -m MYDOMAIN\\Required-AD-Group"

    <Directory "/usr/local/www/apache24/data">
        Order deny,allow
        Allow from all
        AuthType Basic
        AuthName "NTLM auth"
        AuthBasicProvider external
        AuthnCacheProvideFor external
        AuthExternal ntlmAuth
        Require valid-user
    </Directory>
</VirtualHost>


```



## Licensing

"The code in this project is licensed under BSD-new license."