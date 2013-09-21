package main

func AboutHandler(ctx *Context) {
	ctx.Render("about.html", map[string]interface{}{"title": "About"})
}

func HelpHandler(ctx *Context) {
	ctx.Render("help.html", map[string]interface{}{"title": "Help"})
}

func HelpCheatSheetHandler(ctx *Context) {
	ctx.Render("help_cheatsheet.html", map[string]interface{}{"title": "Cheat Sheet"})
}

func HelpGettingStartedHandler(ctx *Context) {
	ctx.Render("help_getting_started.html", map[string]interface{}{"title": "Getting Started"})
}

func HelpGlossaryHandler(ctx *Context) {
	ctx.Render("help_glossary.html", map[string]interface{}{"title": "Glossary"})
}

func HomeHandler(ctx *Context) {
	ctx.Render("home.html", map[string]interface{}{"bodytitle": "Navigate the web faster.", "bodyid": "home"})
}

func PrivacyHandler(ctx *Context) {
	ctx.Render("privacy.html", map[string]interface{}{"title": "Privacy Policy"})
}

func TermsHandler(ctx *Context) {
	ctx.Render("terms.html", map[string]interface{}{"title": "Terms of Service"})
}
