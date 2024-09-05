# Mailonfail

<!-- toc -->

- [Synopsis](#synopsis)
- [Syntax](#syntax)
- [Config](#config)
- [Mail template](#mail-template)

<!-- /toc -->

## Synopsis

Mailtofail sends a notification email when a command fails. I hope that it is useful especially for cron jobs.

## Syntax

Keep in mind that flags have to precede the command that should be executed. Usage examples:

```shell
mailonfail ls -la

# set config file and log level
mailonfail --log-level debug -c myconf.yaml ls -la
```

## Config

An example config looks like the following...

```go mdox-exec="tail -n+2 examples/conf.yaml"
smtp_host: smtp.host.com
smtp_port: 993
smtp_user: test_user
smtp_pass: test_pass
mail_from: user@test.mail
mail_to: user@test.mail
mail_on_success: false

subject_prefix: "[{{.hostname}}] "
mail_template: >
  </br><b>Error</b></br>{{.error}}</br> </br><b>Output</b></br>{{.output}}</br>
  </br><b>Exitcode</b></br></br>{{.exitcode}}</br>
```

Every config value can be overwritten by an env var. This comes in handy in `crontabs`. Available env vars...

```go mdox-exec="sh/print_av_env_vars.sh"
MOF_MAIL_FROM
MOF_MAIL_ON_SUCCESS
MOF_MAIL_TEMPLATE
MOF_MAIL_TO
MOF_SMTP_HOST
MOF_SMTP_PASS
MOF_SMTP_PORT
MOF_SMTP_USER
MOF_SUBJECT_PREFIX
```

If `subject_prefix` and `mail_template` are not set they are loaded from the default [conf](src/default_conf.yaml).

## Mail template

The following variables are available in the mail template. Make sure to use golang template syntax (e.g. `{{.username}}`, `{{.output}}`). Available template vars...

```go mdox-exec="sh/print_av_tpl_vars.sh"
command
error
exitcode
group_id
home
hostname
output
runtime
user
user_id
user_name
```

## Help

```go mdox-exec="r -h"

If a command fails, send a mail...

Usage: mailonfail [flags] [<command> ...]

Arguments:
  [<command> ...]    command to run, flags always have to be in front

Flags:
  -h, --help                      Show context-sensitive help.
  -c, --config-file=STRING        config file to load, values can be overwritten
                                  by env vars
      --log-file="/dev/stdout"    log file
      --log-level="info"          log level
      --log-no-colors             disable output colours, print plain text
      --log-json                  enable json log, instead of text one
  -n, --dry-run                   dry run, just print operations that would run
  -V, --version-flag              display version
```
