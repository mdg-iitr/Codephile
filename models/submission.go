package models

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scripts"
	"log"
)

func AddSubmissions(user *types.User, site string) error {
	var handle string
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	switch site {
	case CODECHEF:
		handle = user.Handle.Codechef
		if handle == "" {
			return errors.New("handle not available")
		}
		addSubmissions := scripts.GetCodechefSubmissions(handle, user.Last.Codechef)
		if len(addSubmissions) != 0 {
			user.Last.Codechef = addSubmissions[0].CreationDate
			change := bson.M{"$push": bson.M{"submission.codechef": bson.M{"$each": addSubmissions}}, "$set": bson.M{"lastfetched": user.Last}}
			err := coll.Collection.UpdateId(user.ID, change)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		return nil
	case CODEFORCES:
		handle = user.Handle.Codeforces
		if handle == "" {
			return errors.New("handle not available")
		}
		fmt.Println(user.Last.Codeforces)
		addSubmissions := scripts.GetCodeforcesSubmissions(handle, user.Last.Codeforces).Data
		if len(addSubmissions) != 0 {
			user.Last.Codeforces = addSubmissions[0].CreationDate
			change := bson.M{"$push": bson.M{"submission.codeforces": bson.M{"$each": addSubmissions}}, "$set": bson.M{"lastfetched": user.Last}}
			err := coll.Collection.UpdateId(user.ID, change)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		return nil
	case SPOJ:
		handle = user.Handle.Spoj
		if handle == "" {
			return errors.New("handle not available")
		}
		addSubmissions := scripts.GetSpojSubmissions(handle, user.Last.Spoj)
		if len(addSubmissions) != 0 {
			user.Last.Spoj = addSubmissions[0].CreationDate
			change := bson.M{"$push": bson.M{"submission.spoj": bson.M{"$each": addSubmissions}}, "$set": bson.M{"lastfetched": user.Last}}
			err := coll.Collection.UpdateId(user.ID, change)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		return nil
	case HACKERRANK:
		handle = user.Handle.Hackerrank
		if handle == "" {
			return errors.New("handle not available")
		}
		addSubmissions := scripts.GetHackerrankSubmissions(handle, user.Last.Hackerrank).Data
		if len(addSubmissions) != 0 {
			user.Last.Hackerrank = addSubmissions[0].CreationDate
			change := bson.M{"$push": bson.M{"submission.hackerrank": bson.M{"$each": addSubmissions}}, "$set": bson.M{"lastfetched": user.Last}}
			err := coll.Collection.UpdateId(user.ID, change)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		return nil
	}
	return nil
}

func GetSubmissions(ID bson.ObjectId) (*types.Submissions, error) {
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	var user types.User
	err := coll.Collection.FindId(ID).Select(bson.M{"submission": 1}).One(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user.Submissions, nil
}

func FilterSubmission(uid bson.ObjectId, status string, tag string, site string) ([]map[string]interface{}, error) {
	c := db.NewUserCollectionSession()
	defer c.Close()
	fmt.Println(status)
	match1 := bson.M{
		"$match": bson.M{
			"_id": uid,
		},
	}
	unwind := bson.M{
		"$unwind": "$submission." + site,
	}
	match2 := bson.M{
		"$match": bson.M{"submission." + site + ".status": status},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":                0,
			"submission." + site: 1,
		},
	}
	all := []bson.M{match1, unwind, match2, project}
	pipe := c.Collection.Pipe(all)

	var result map[string]interface{}
	iter := pipe.Iter()
	var final []map[string]interface{}
	for iter.Next(&result) {
		final = append(final, result["submission"].(map[string]interface{})[site].(map[string]interface{}))
	}
	return final, nil
}
