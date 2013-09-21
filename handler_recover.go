package main

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
	"net/smtp"
)

func generateRecoverEmail(user *User) error {
	token := uuid4()
	if err := Config.User.Update(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"recover": token}}); err == nil {
		auth := smtp.PlainAuth("", Config.EmailAddress, Config.EmailPassword, Config.SmtpServer)
		body := fmt.Sprintf(
			"to: %s\nsubject: KeyFu Password Reset Request\n"+
				"A password request has been submitted to KeyFu for your account. If this email "+
				"has reached you in error, you can safely ignore it.\n\n"+
				"Account: %s\n\n"+
				"To reset your password, follow the link provided here:\n\n"+
				"https://www.keyfu.com/recover?token=%s", user.Email, user.Email, token)
		return smtp.SendMail(Config.SmtpServer+":25", auth, Config.EmailAddress, []string{user.Email}, []byte(body))
	}

	return errors.New("Failed to set recover token.")
}

func getUserByResetToken(token string) *User {
	user := User{}
	err := Config.User.Find(bson.M{"recover": token}).One(&user)
	if err == nil {
		return &user
	}

	return nil
}

func RecoverHandler(ctx *Context) {
	values := ctx.Request.URL.Query()
	if token := values.Get("token"); token != "" {
		if user := getUserByResetToken(token); user != nil {
			ctx.Render("recover_reset.html", map[string]interface{}{"title": "Reset Password", "token": token})
		} else {
			ctx.SetError("Password reset token expired.")
			ctx.Redirect("/")
		}
	} else {
		ctx.Render("recover.html", map[string]interface{}{"title": "Recover Account"})
	}
}

func RecoverPostHandler(ctx *Context) {
	if token := ctx.Request.FormValue("token"); token != "" {
		password := ctx.Request.FormValue("new-password")

		content := map[string]interface{}{
			"title": "Reset Password",
			"token": token,
		}

		if user := getUserByResetToken(token); user != nil {
			if err := user.SetPassword(password); err == nil {
				if err = Config.User.Update(bson.M{"_id": user.Id},
					bson.M{"$set": bson.M{"password": user.Password}, "$unset": bson.M{"recover": 1}}); err == nil {
					ctx.SetSuccess("Password successfully reset.")
					ctx.Redirect("/login")
					return
				} else {
					content["error"] = "Failed to update password."
				}
			} else {
				content["error"] = err.Error()
			}
		} else {
			ctx.SetError("Password reset token expired.")
			ctx.Redirect("/")
			return
		}

		ctx.Render("recover_reset.html", content)
	} else {
		email := ctx.Request.FormValue("email")

		content := map[string]interface{}{
			"title": "Recover Account",
			"email": email,
		}

		user := User{}
		err := Config.User.Find(bson.M{"email": email}).One(&user)
		if err == nil {
			if err = generateRecoverEmail(&user); err == nil {
				ctx.SetSuccess("Password reset email sent.")
				ctx.Redirect("/")
				return
			} else {
				content["error"] = "Failed to send email, try again later."
			}
		} else {
			content["error"] = "Email address not found."
		}

		ctx.Render("recover.html", content)
	}
}
