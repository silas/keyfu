package main

import (
	"html/template"
	"labix.org/v2/mgo"
)

type pathRoute map[string]Handler

type ServerConfig struct {
	EmailAddress        string
	EmailPassword       string
	SmtpServer          string
	SessionLength       int64
	SessionCookieSecret string
	MongoHost           string
	Templates           *template.Template
	CdnUrl              string
	CssPath             string
	JsPath              string
	JsBin               string
	LogoPath            string
	OpenSearchPath      string
	Catalog             *mgo.Collection
	Keyword             *mgo.Collection
	User                *mgo.Collection
}

var Config = &ServerConfig{
	EmailAddress:        "help@keyfu.com",
	EmailPassword:       "xx7<F.YoR%M=<\"]KiC(9*?$UsQF=-QIk",
	SmtpServer:          "smtp.gmail.com",
	SessionLength:       30 * 24 * 60 * 60,
	SessionCookieSecret: "JpxjF_e^}3-[uVn'v:^G/,I.RuaSk 0&",
	MongoHost:           "localhost",
	CssPath:             "/static/keyfu.css",
	JsPath:              "/static/keyfu.js",
	LogoPath:            "/static/logo.gif",
	OpenSearchPath:      "/static/opensearch.xml",
	CdnUrl:              "http://c658141.r41.cf2.rackcdn.com",
}

func run(httpInterface string, staticPath string, templatesPath string) {
	if t, err := template.ParseGlob(templatesPath + "/*.html"); err == nil {
		Config.Templates = t
	} else {
		panic(err)
	}

	Config.CssPath = versionStatic(staticPath, Config.CssPath)
	Config.JsPath = versionStatic(staticPath, Config.JsPath)
	Config.LogoPath = versionStatic(staticPath, Config.LogoPath)
	Config.OpenSearchPath = versionStatic(staticPath, Config.OpenSearchPath)

	mongoConnection, err := mgo.Dial(Config.MongoHost)
	if err != nil {
		panic(err)
	}

	defer mongoConnection.Close()
	mongoConnection.SetMode(mgo.Monotonic, true)

	db := mongoConnection.DB("keyfu")
	Config.Catalog = db.C("catalog")
	Config.Keyword = db.C("keyword")
	Config.User = db.C("user")

	catalogIndex := mgo.Index{
		Key:        []string{"value"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	err = Config.Catalog.EnsureIndex(catalogIndex)
	if err != nil {
		panic(err)
	}

	keywordIndex := mgo.Index{
		Key:        []string{"uid", "key"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = Config.Keyword.EnsureIndex(keywordIndex)
	if err != nil {
		panic(err)
	}

	userIndex := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = Config.User.EnsureIndex(userIndex)
	if err != nil {
		panic(err)
	}

	// Setup hardcoded shortcuts
	builtinKeywords["account"] = builtinKeywords["com.keyfu.account"]
	builtinKeywords["a"] = builtinKeywords["com.keyfu.add"]
	builtinKeywords["add"] = builtinKeywords["com.keyfu.add"]
	builtinKeywords["cheatsheet"] = builtinKeywords["com.keyfu.help.cheatsheet"]
	builtinKeywords["cs"] = builtinKeywords["com.keyfu.help.cheatsheet"]
	builtinKeywords["d"] = builtinKeywords["com.keyfu.delete"]
	builtinKeywords["delete"] = builtinKeywords["com.keyfu.delete"]
	builtinKeywords["e"] = builtinKeywords["com.keyfu.edit"]
	builtinKeywords["edit"] = builtinKeywords["com.keyfu.edit"]
	builtinKeywords["l"] = builtinKeywords["com.keyfu.list"]
	builtinKeywords["list"] = builtinKeywords["com.keyfu.list"]
	builtinKeywords["h"] = builtinKeywords["com.keyfu.help"]
	builtinKeywords["help"] = builtinKeywords["com.keyfu.help"]
	builtinKeywords["home"] = builtinKeywords["com.keyfu"]
	builtinKeywords["login"] = builtinKeywords["com.keyfu.login"]
	builtinKeywords["logout"] = builtinKeywords["com.keyfu.logout"]
	builtinKeywords["privacy"] = builtinKeywords["com.keyfu.privacy"]
	builtinKeywords["r"] = builtinKeywords["com.keyfu.rename"]
	builtinKeywords["rename"] = builtinKeywords["com.keyfu.rename"]
	builtinKeywords["terms"] = builtinKeywords["com.keyfu.terms"]

	s := NewServer()
	s.Get("/run", RunHandler)
	s.Get("/autocomplete", AutocompleteHandler)
	s.Get("/", HomeHandler)
	s.Get("/add", LoginHelper(AddHandler))
	s.Get("/about", AboutHandler)
	s.Get("/account", LoginHelper(AccountHandler))
	s.Post("/account", LoginCsrfHelper(AccountPostHandler))
	s.Get("/delete", LoginHelper(DeleteHandler))
	s.Post("/delete", LoginCsrfHelper(DeletePostHandler))
	s.Get("/edit", LoginHelper(EditHandler))
	s.Post("/edit", LoginCsrfHelper(EditPostHandler))
	s.Get("/help", HelpHandler)
	s.Get("/help/cheatsheet", HelpCheatSheetHandler)
	s.Get("/help/getting-started", HelpGettingStartedHandler)
	s.Get("/help/glossary", HelpGlossaryHandler)
	s.Get("/list", LoginHelper(ListHandler))
	s.Get("/login", LoginHandler)
	s.Post("/login", CsrfHelper(LoginPostHandler))
	s.Get("/logout", LogoutHandler)
	s.Get("/privacy", PrivacyHandler)
	s.Get("/recover", RecoverHandler)
	s.Post("/recover", CsrfHelper(RecoverPostHandler))
	s.Get("/rename", LoginHelper(RenameHandler))
	s.Post("/rename", LoginCsrfHelper(RenamePostHandler))
	s.Get("/signup", SignupHandler)
	s.Post("/signup", CsrfHelper(SignupPostHandler))
	s.Get("/terms", TermsHandler)
	s.Run(httpInterface, staticPath)
}
