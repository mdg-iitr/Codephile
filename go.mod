module github.com/mdg-iitr/Codephile

// +heroku install ./cmd/... .

go 1.12

require (
	cloud.google.com/go/firestore v1.0.0 // indirect
	cloud.google.com/go/storage v1.10.0
	firebase.google.com/go v3.9.0+incompatible
	github.com/PuerkitoBio/goquery v1.5.0 // indirect
	github.com/antchfx/htmlquery v1.1.0 // indirect
	github.com/antchfx/xmlquery v1.1.0 // indirect
	github.com/antchfx/xpath v1.1.0 // indirect
	github.com/astaxie/beego v1.12.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/sentry-go v0.5.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gocolly/colly v1.2.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/schema v1.1.0
	github.com/joho/godotenv v1.3.0
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/olivere/elastic/v7 v7.0.6
	github.com/pkg/errors v0.8.1
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/temoto/robotstxt v1.1.1 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78
	google.golang.org/api v0.30.0
)
