// Copyright (c) 2017 Kuguar <licenses@kuguar.io> Author: Adrian P.K. <apk@kuguar.io>
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package tests

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/testbootstrap"

	_ "github.com/lib/pq"
)

var (
	tbp         = testbootstrap.TestBootstrap
	usersURL    string
	profilesURL string
	avatarURL   string
	user1       = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2       = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	profile1    = "28bb0dad-ece8-44a1-8c45-c4898968bee5"
	profile2    = "1b57cb73-7f61-4323-ae87-86b4d0569178"
)

func init() {
	usersURL = fmt.Sprintf("%s/users", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

// func TestGetProfileByUserId(t *testing.T) {
// 	logger.Debug("TestGetProfileByUserID...")
// 	prepareTestDatabase()
// 	reader = strings.NewReader("")
// 	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, "5958b185-8150-4aae-b53f-0c44771ddec5")
// 	request, _ := http.NewRequest("GET", profileURL, reader)
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Errorf("Error executing request: %s", err.Error())
// 		return
// 	}
// 	if res.StatusCode != http.StatusOK {
// 		t.Errorf("Status: %d | Expected: 200-StatusOk", res.StatusCode)
// 	}
// }

func TestPrepare(t *testing.T) {
	logger.Debug("TestUploadAvatar...")
	tbp.PrepareTestDatabase()
}

func TestUploadAvatar(t *testing.T) {
	logger.Debug("TestUploadAvatar...")
	tbp.PrepareTestDatabase()
	base64Image := sampleBase64Image()
	avatarJSON := fmt.Sprintf(`
  {
    "data": {
      "profileID": "%s",
      "base64": "%s"
    }
  }
	`, profile1, base64Image)
	//data:image/png;base64,
	//logger.Debugf("Sample Base64 encoded image in json: '%s'", avatarJSON)
	tbp.Reader = strings.NewReader(avatarJSON)
	avatarURL := fmt.Sprintf("%s/%s/profile/avatar", usersURL, user1)
	request, _ := http.NewRequest("POST", avatarURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Error(err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

// func TestUpdateProfile(t *testing.T) {
// 	logger.Debug("TestUpdateProfile...")
// 	prepareTestDatabase()
// 	profileJSON := `
// 	{
// 		"data": {
// 			"id": 1,
// 			"name": "Admin",
// 			"email": "admin@gmail.com",
// 		  "description": "But you can callme Admy",
// 		  "bio": "I like iron.",
// 		  "moto": "I do my job",
// 		  "website": "admin.com",
// 		  "anniversary_date": 1420070400,
// 			"user_id": 1
// 		}
// 	}
// 	`
// 	reader = strings.NewReader(profileJSON)
// 	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, "1")
// 	request, _ := http.NewRequest("PUT", profileURL, reader)
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Errorf("Error executing request: %s", err.Error())
// 		return
// 	}
// 	if res.StatusCode != http.StatusNoContent {
// 		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
// 	}
// }
//
// func TestUpdateProfileWithDatabaseVerify(t *testing.T) {
// 	logger.Debug("TestUpdateUserWithDatabaseVerify...")
// 	prepareTestDatabase()
// 	newName := "Admin"
// 	newEmail := "admin@gmail.com"
// 	profileJSON := fmt.Sprintf(`
// 	{
// 		"data": {
// 			"id": 1,
// 			"name": "%s",
// 			"email": "%s",
// 			"description": "But you can callme Admy",
// 			"bio": "I like iron.",
// 			"moto": "I do my job",
// 			"website": "admin.com",
// 			"anniversary_date": 1420070400,
// 			"user_id": 1
// 		}
// 	}
// 	`, newName, newEmail)
// 	reader = strings.NewReader(profileJSON)
// 	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, "1")
// 	request, _ := http.NewRequest("PUT", profileURL, reader)
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Errorf("Error executing request: %s", err.Error())
// 		return
// 	}
// 	if res.StatusCode == http.StatusNoContent {
// 		profileRepo, err := repo.MakeProfileRepository()
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		profile, err := profileRepo.GetByUserIDStr("1")
// 		if err == nil {
// 			if profile.Name.String == newName && profile.Email.String == newEmail {
// 				logger.Debug("Profile update: ok.")
// 			} else {
// 				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", profile.Name.String, newName)
// 				error += fmt.Sprintf("Email: '%s' | Expected: '%s'", profile.Email, newEmail)
// 				t.Error(error)
// 			}
// 		} else {
// 			t.Error(err.Error())
// 		}
// 	} else {
// 		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
// 	}
// }
//
// func TestDeleteProfile(t *testing.T) {
// 	logger.Debug("TestDeleteProfile...")
// 	prepareTestDatabase()
// 	reader = strings.NewReader("")
// 	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, "1")
// 	request, _ := http.NewRequest("DELETE", profileURL, reader)
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Errorf("Error executing request: %s", err.Error())
// 		return
// 	}
// 	if res.StatusCode != http.StatusNoContent {
// 		t.Errorf("Status: %d | Expected: 200-StatusNoContent", res.StatusCode)
// 	}
// }
//
// func TestDeleteProfileWithDatabaseVerify(t *testing.T) {
// 	logger.Debug("TestDeleteProfileWithDatabaseVerify...")
// 	prepareTestDatabase()
// 	reader = strings.NewReader("")
// 	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, "1")
// 	request, _ := http.NewRequest("DELETE", profileURL, reader)
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Errorf("Error executing request: %s", err.Error())
// 		return
// 	}
// 	if res.StatusCode == http.StatusNoContent {
// 		profileRepo, err := repo.MakeProfileRepository()
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		profile, err := profileRepo.GetByUserIDStr("1")
// 		if err != nil {
// 			logger.Debug("TestDeleteProfile: ok")
// 		} else {
// 			t.Errorf("Profile: %s | Expected: 'nil'", profile.Name)
// 		}
// 	} else {
// 		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
// 	}
// }

func sampleBase64Image() string {
	pngFilePath := currentDir()
	pngFile := path.Join("resources", "testdata", "sampleavatar.png")
	imgFile, err := os.Open(path.Join(pngFilePath, pngFile))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer imgFile.Close()
	// Create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buffer := make([]byte, size)
	// Read file content into buffer
	reader := bufio.NewReader(imgFile)
	reader.Read(buffer)
	// If you create a new image instead of loading from file, encode the image to buffer instead with png.Encode()
	// png.Encode(&buf, image)
	// Convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buffer)
	return imgBase64Str
}

func currentDir() string {
	exPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	parentPath := filepath.Dir(exPath)
	return parentPath
}
