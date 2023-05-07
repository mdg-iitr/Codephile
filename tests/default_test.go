package test

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"

	"github.com/astaxie/beego"
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/types"
	_ "github.com/mdg-iitr/Codephile/routers"
	"github.com/mdg-iitr/Codephile/services/auth"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	beego.TestBeegoInit(conf.AppRootDir)
	db.NewUserCollectionSession().DropDatabase()
}

// TestGet is a sample to run an endpoint test
func TestGetAllUsers(t *testing.T) {
	uid, _ := models.AddUser(types.User{
		ID:        primitive.NewObjectID(),
		Email:     "test@abc.com",
		Username:  "test",
		FullName:  "Test User",
		Institute: "IIT Roorkee",
		Password:  "password",
	})
	token := auth.GenerateToken(uid)
	r, _ := http.NewRequest("GET", "/v1/user/all", nil)
	r.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}
