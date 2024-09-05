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
	subject := conf.execTemplate(conf.SubjectPrefix, cr)
	body := conf.execTemplate(conf.MailTemplate, cr)
	lg.Trace("send mail", logseal.F{"body": body, "subject": subject})
	if !conf.DryRun {
		m := mail.NewMessage()
		m.SetHeader("From", conf.MailFrom)
		m.SetHeader("To", conf.MailTo)
		m.SetHeader("Subject", subject)
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
	ui := getUserInfo()
	templ := template.Must(template.New("tpl").Parse(s))
	err := templ.Execute(buf, map[string]interface{}{
		"cmd":       fmt.Sprintf("%q", conf.Cmd),
		"runtime":   cr.RunTime,
		"error":     cr.Error,
		"exitcode":  cr.Exitcode,
		"output":    cr.Out,
		"hostname":  getHostName(),
		"user_id":   ui.UserID,
		"group_id":  ui.GroupID,
		"user":      ui.UserName,
		"user_name": ui.Name,
		"home":      ui.Home,
	})
	lg.IfErrError("unable to use mail template", logseal.F{"error": err})
	return buf.String()
}

func getHostName() (hostname string) {
	hostname, _ = os.Hostname()
	return
}

type userInfo struct {
	UserID   string
	GroupID  string
	UserName string
	Name     string
	Home     string
}

func getUserInfo() (ui userInfo) {
	user, err := user.Current()
	if err == nil {
		ui.UserID = user.Uid
		ui.GroupID = user.Gid
		ui.UserName = user.Username
		ui.Name = user.Name
		ui.Home = user.HomeDir
	}
	return ui
}
