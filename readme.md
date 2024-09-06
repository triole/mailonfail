# Mailonfail

<!-- toc -->

- [Synopsis](#synopsis)
- [Syntax](#syntax)
- [Config](#config)
- [Mail template](#mail-template)

<!-- /toc -->

## Synopsis

Mailtofail sends a notification email when a command fails. It might be useful especially for cron jobs.

## Syntax

Keep in mind that flags have to precede the command that should be executed. Usage examples:

```shell
mailonfail ls -la

# set config file and log level, and run "ls -la"
mailonfail --log-level debug -c myconf.yaml ls -la
```

## Config

An example config looks like the following...

```go mdox-exec="tail -n+2 examples/conf.yaml"
smtp_host: smtp.host.com
smtp_port: 2525
smtp_user: test_user
smtp_pass: test_pass
mail_from: user@test.mail
mail_to: user@test.mail
mail_on_success: false

mail_subject: >
  [mof@{{.hostname}}] ({{if .success}}success{{else}}error{{end}})
  cmd {{.command}}, exitcode {{.exitcode}}
mail_body: |
  <b>Run Start</b>&nbsp;&nbsp;{{.run_start}}</br>
  <b>Run End</b>&nbsp;&nbsp;&nbsp;&nbsp;{{.run_end}}</br>
  <b>Duration</b>&nbsp;&nbsp;&nbsp;{{.run_duration}}</br></br>
  <b>User</b>&nbsp;&nbsp;{{.user}}</br>
  <b>Command</b>&nbsp;&nbsp;{{.command}}</br>
  </br><b>Output</b></br><pre>{{.output}}</pre></br>
  {{if .error}}
  </br><b>Error</b></br>{{.error}}</br>
  {{end}}
  </br><b>Exitcode</b></br>{{.exitcode}}</br>
```

Every config value can be overwritten by an env var. This comes in handy in `crontabs`. Available env vars...

```go mdox-exec="sh/print_av_env_vars.sh"
MOF_SMTP_HOST
MOF_SMTP_PORT
MOF_SMTP_USER
MOF_SMTP_PASS
MOF_MAIL_FROM
MOF_MAIL_TO
MOF_MAIL_ON_SUCCESS
MOF_MAIL_SUBJECT
MOF_MAIL_BODY
```

If `mail_subject` or `mail_body` are not set, they are loaded from the default [conf](src/default_conf.yaml).

## Mail template

The following variables are available in the email templates. Make sure to use golang template syntax (e.g. `{{.username}}`, `{{.output}}`).

```go mdox-exec="sh/print_av_tpl_vars.sh"
run_start
run_end
run_duration
command
output
error
exitcode
success
hostname
user_id
group_id
user
user_name
home
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
