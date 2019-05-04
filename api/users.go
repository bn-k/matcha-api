package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserModify(c *gin.Context) {
	var req Request
	if err := req.prepareRequest(c); err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.body = getBodyToMap(c)
		req.user, _ = app.getUser(req.id, "")
		mod := c.Param("name")
		switch mod {
		case "position":
			req.updatePosition()
			break
		case "location":
			req.updateLocation()
			break
		case "biography":
			req.updateBio()
			break
		case "username":
			req.updateUsername()
			break
		case "tag":
			req.addTag()
			break
		case "password":
			req.updatePassword()
			break
		case "firstname":
			req.updateFirstname()
			break
		case "lastname":
			req.updateLastname()
			break
		case "genre":
			req.updateGenre()
			break
		case "email":
			req.updateEmail()
			break
		case "interest":
			req.updateInterest()
			break
		default:
			c.JSON(202, gin.H{"err": "route not found"})
		}
	}
}

func (req Request) updatePosition() {
	pos := req.body["position"]
	fmt.Println(pos)
	retUser(req)
}
func (req Request) updateLocation() {
	pos := req.body["location"]
	fmt.Println(pos)
	retUser(req)
}
func (req Request) updateBio() {
	bio := req.body["biography"].(string)

	if len(bio) > 100 || len(bio) < 10 {
		err := errors.New("error : your biography must be between 10 and 100 characters")
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.user.Biography = bio
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) checkPassword() error {
	pass := req.body["old_password"].(string)
	truePass := hash.Decrypt(hashKey, req.user.Password)
	if pass != truePass {
		return errors.New("Wrong password")
	} else {
		return nil
	}
}

func (req Request) updateUsername() {
	newUsername := req.body["new_username"].(string)
	if err := req.checkPassword(); err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else if len(newUsername) < 6 || len(newUsername) > 20 {
		err = errors.New("error : your username must be between 6 to 20 characters")
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.user.Username = newUsername
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) updateGenre() {
	genre := req.body["genre"].(string)
	req.user.Genre = genre
	app.updateUser(req.user)
	retUser(req)
}

func (req Request) updateInterest() {
	interest := req.body["interest"].(string)
	req.user.Interest = interest
	app.updateUser(req.user)
	retUser(req)
}

func (req Request) addTag() {
}

func (req Request) updateEmail() {
	newEmail := req.body["new_email"].(string)
	if err := req.checkPassword(); err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else if !emailIsValid(newEmail) {
		req.context.JSON(201, gin.H{"err": "Email is invalid"})
	} else {
		req.user.Email = newEmail
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) updatePassword() {
	newPassword := req.body["new_password"].(string)
	confirmPassword := req.body["confirm"].(string)

	if err := req.checkPassword(); err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		if err = verifyPassword(newPassword, confirmPassword); err != nil {
			req.context.JSON(201, gin.H{"err": err.Error()})
		} else {
			req.user.Password = hash.Encrypt(hashKey, newPassword)
			app.updateUser(req.user)
			retUser(req)
		}
	}
}

func (req Request) updateFirstname() {
	firstname := req.body["firstname"].(string)
	if len(firstname) < 2 || len(firstname) > 20 {
		req.context.JSON(201, gin.H{"err": "error : your firstname must be between 2 to 20 characters"})
	} else {
		req.user.FirstName = firstname
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) updateLastname() {
	lastname := req.body["lastname"].(string)
	if len(lastname) < 2 || len(lastname) > 20 {
		err := errors.New("error : your lastname must be between 2 to 20 characters")
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.user.LastName = lastname
		app.updateUser(req.user)
		retUser(req)
	}
}

func UserImageHandler(c *gin.Context) {
	file := c.PostForm("file")
	fmt.Printf("file  %s\n", file)

	claims := jwt.MapClaims{}
	valid, err := ValidateToken(c, &claims)

	if valid == true {
		Id := int(claims["id"].(float64))
		g, err := app.dbGetUserProfile(Id)
		tagList := app.dbGetTagList()
		if err != nil {
			c.JSON(201, gin.H{"err": err.Error()})
		} else {
			c.JSON(200, gin.H{"user": g, "tagList": tagList})
		}
	} else {
		c.JSON(201, gin.H{"err": err.Error()})
	}
}

type Request struct {
	context *gin.Context
	claims  jwt.MapClaims
	user    User
	body    map[string]interface{}
	id      int
}

func getBodyToMap(c *gin.Context) (body map[string]interface{}) {
	r, _ := c.GetRawData()
	err := json.Unmarshal(r, &body)
	if err != nil {
		panic(err)
	}
	return
}

func (req *Request) prepareRequest(c *gin.Context) error {
	req.context = c
	req.claims = jwt.MapClaims{}
	valid, err := ValidateToken(c, &req.claims)
	if valid == true {
		req.id = int(req.claims["id"].(float64))
		req.user, _ = app.getUser(req.id, "")
	} else {
		return err
	}
	return nil
}

func UserHandler(c *gin.Context) {
	var req Request
	if err := req.prepareRequest(c); err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
	}
	retUser(req)
}

func retUser(req Request) {
	g, err := app.dbGetUserProfile(req.id)
	tagList := app.dbGetTagList()
	userTags := app.dbGetUserTags(req.user.Username)
	if err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.context.JSON(200, gin.H{"user": g, "tagList": tagList, "userTags": userTags})
	}
}
