package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"text/template"

	"github.com/triole/logseal"
	mail "gopkg.in/mail.v2"
)

func (conf tConf) sendMail(cr cmdReturn) {
	body := conf.execTemplate(conf.MailTemplate, cr)
	lg.Trace("send mail", logseal.F{"body": body})
	if !conf.DryRun {
		m := mail.NewMessage()
		m.SetHeader("From", conf.MailFrom)
		m.SetHeader("To", conf.MailTo)
		m.SetHeader("Subject", conf.execTemplate(conf.SubjectPrefix, cr)+"cmd: "+fmt.Sprintf("%s", conf.Cmd))
		m.SetBody("text/html", body)
		d := mail.NewDialer(
			conf.SmtpHost, conf.SmtpPort, conf.SmtpUser, conf.SmtpPass,
		)
		if err := d.DialAndSend(m); err != nil {
			lg.IfErrError("can not send mail", logseal.F{"error": err})
		}
	}
}

func (conf tConf) execTemplate(s string, cr cmdReturn) string {
	buf := &bytes.Buffer{}
	templ := template.Must(template.New("tpl").Parse(s))
	err := templ.Execute(buf, map[string]interface{}{
		"command":  conf.Cmd,
		"runtime":  cr.RunTime,
		"error":    cr.Error,
		"exitcode": cr.Exitcode,
		"output":   cr.Out,
		"hostname": getHostName(),
		"user":     getUserName(),
	})
	lg.IfErrError("unable to use mail template", logseal.F{"error": err})
	return buf.String()
}

func getHostName() (hostname string) {
	hostname, _ = os.Hostname()
	return
}

func getUserName() (userName string) {
	user, err := user.Current()
	if err == nil {
		userName = user.Username
	}
	return
}
