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

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "image/gif"  // Required by "image" library.
	_ "image/jpeg" // Required by "image" library.
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/adrianpk/fundacja/app"
	"github.com/adrianpk/fundacja/models"
	"github.com/adrianpk/fundacja/repo"
)

const (
	website     = "https://blueimp.github.io/jQuery-File-Upload/"
	minFileSize = 1 // bytes
	// Max file size is memcache limit (1MB) minus key size minus overhead:
	maxFileSize            = 999000 // bytes
	imageTypesPattern      = "image/(gif|p?jpeg|(x-)?png)"
	acceptFileTypesPattern = imageTypesPattern
	thumbMaxWidth          = 80
	thumbMaxHeight         = 80
	expirationTime         = 300 // seconds
	// If empty, only allow redirects to the referer protocol+host.
	// Set to a regexp string for custom pattern matching:
	redirectAllowTargetPattern = ""
)

var (
	imageTypes      = regexp.MustCompile(imageTypesPattern)
	acceptFileTypes = regexp.MustCompile(acceptFileTypesPattern)
	thumbSuffix     = "." + fmt.Sprint(thumbMaxWidth) + "x" + fmt.Sprint(thumbMaxHeight)
)

// HandleAvatar - Handles avatar REST related functions
func HandleAvatar(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		app.ShowError(w, app.ErrRequest, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add(
		"Access-Control-Allow-Methods",
		"OPTIONS, HEAD, GET, POST, DELETE",
	)
	w.Header().Add(
		"Access-Control-Allow-Headers",
		"Content-Type, Content-Range, Content-Disposition",
	)
	switch r.Method {
	case "OPTIONS", "HEAD":
		return
	case "GET":
		//get(w, r)
	case "POST":
		if len(params["_method"]) > 0 && params["_method"][0] == "DELETE" {
			//delete(w, r)
		} else {
			post(w, r)
		}
	case "DELETE":
		//delete(w, r)
	default:
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	result := make(map[string][]*models.Image, 1)
	result["images"] = handleJSONUploads(w, r)
	b, err := json.Marshal(result)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	if redirect := r.FormValue("redirect"); validateRedirect(r, redirect) {
		if strings.Contains(redirect, "%s") {
			redirect = fmt.Sprintf(
				redirect,
				escape(string(b)),
			)
		}
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	jsonType := "application/json"
	if strings.Index(r.Header.Get("Accept"), jsonType) != -1 {
		w.Header().Set("Content-Type", jsonType)
	}
}

func handleJSONUploads(w http.ResponseWriter, r *http.Request) (fileInfos []*models.Image) {
	fileInfos = make([]*models.Image, 0)
	var res AvatarResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	base64Data := res.Data.Base64
	// Associated Profile
	profile, err := profileFromSessionID(r)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Get repo
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist

	profile.SetAvatarAsBase64(base64Data)
	profileRepo.SaveProfileAvatar(&profile)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Output
	w.WriteHeader(http.StatusNoContent)
	return
}

// func storeAsFileOut(imageBin []byte) {
// 	logger.Debug("Storing image as a file.")
// 	//imgByte := code.PNG()
// 	// convert []byte to image for saving to file
// 	img, err2 := png.Decode(bytes.NewReader(imageBin))
// 	logger.Debugf("Error 2 %v", err2)
// 	logger.Debugf("*** La imagen pasada en bytes es %v", img)
// 	//save the imgByte to file
// 	out, err := os.Create("~/tmp/outcomming.png")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	err = png.Encode(out, img)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

// func storeAsFile(imgByteArray []byte) {
// 	logger.Debug("Storing image as a file.")
// 	//imgByte := code.PNG()
// 	// convert []byte to image for saving to file
// 	img, _, _ := image.Decode(bytes.NewReader(imgByteArray))
// 	//save the imgByte to file
// 	out, err := os.Create("~/tmp/incomming.png")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	err = png.Encode(out, img)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

func escape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

func extractKey(r *http.Request) string {
	// Use RequestURI instead of r.URL.Path, as we need the encoded form:
	path := strings.Split(r.RequestURI, "?")[0]
	// Also adjust double encoded slashes:
	return strings.Replace(path[1:], "%252F", "%2F", -1)
}

// FileInfo - Stores uploaded files related info.
type FileInfo struct {
	Key string `json:"-"`
	//ThumbnailKey string `json:"-"`
	//Url          string `json:"url,omitempty"`
	//ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	//Base64
	//ThumbnailBase6
	Name  string `json:"name"`
	Type  string `json:"type"`
	Size  int64  `json:"size"`
	Error string `json:"error,omitempty"`
	//DeleteUrl    string `json:"deleteUrl,omitempty"`
	//DeleteType   string `json:"deleteType,omitempty"`
}

func getFormValue(p *multipart.Part) string {
	var b bytes.Buffer
	io.CopyN(&b, p, int64(1<<20)) // Copy max: 1 MiB
	return b.String()
}

// image := new ImageResource
// image.ID              bson.ObjectId `bson:"_id,omitempty" json:"id"`
// image.Name            string        `json:"name"`
// image.Description     string        `json:"description"`
// image.Base64          string        `json:"base-64"`
// image.ThumbnailBase64 string        `json:"thumbnail-base-64"`
// image.Location        Geolocation   `bson:"location, omitempty"`
// image.Set

func validateRedirect(r *http.Request, redirect string) bool {
	if redirect != "" {
		var redirectAllowTarget *regexp.Regexp
		if redirectAllowTargetPattern != "" {
			redirectAllowTarget = regexp.MustCompile(redirectAllowTargetPattern)
		} else {
			referer := r.Referer()
			if referer == "" {
				return false
			}
			refererURL, err := url.Parse(referer)
			if err != nil {
				return false
			}
			redirectAllowTarget = regexp.MustCompile("^" + regexp.QuoteMeta(
				refererURL.Scheme+"://"+refererURL.Host+"/",
			))
		}
		return redirectAllowTarget.MatchString(redirect)
	}
	return false
}

// func get(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path == "/" {
// 		http.Redirect(w, r, WEBSITE, http.StatusFound)
// 		return
// 	}
// 	// Use RequestURI instead of r.URL.Path, as we need the encoded form:
// 	key := extractKey(r)
// 	parts := strings.Split(key, "/")
// 	if len(parts) == 3 {
// 		context := appengine.NewContext(r)
// 		item, err := memcache.Get(context, key)
// 		if err == nil {
// 			w.Header().Add("X-Content-Type-Options", "nosniff")
// 			contentType, _ := url.QueryUnescape(parts[0])
// 			if !imageTypes.MatchString(contentType) {
// 				contentType = "application/octet-stream"
// 			}
// 			w.Header().Add("Content-Type", contentType)
// 			w.Header().Add(
// 				"Cache-Control",
// 				fmt.Sprintf("public,max-age=%d", EXPIRATION_TIME),
// 			)
// 			w.Write(item.Value)
// 			return
// 		}
// 	}
// 	http.Error(w, "404 Not Found", http.StatusNotFound)
// }
//
// func delete(w http.ResponseWriter, r *http.Request) {
// 	key := extractKey(r)
// 	parts := strings.Split(key, "/")
// 	if len(parts) == 3 {
// 		result := make(map[string]bool, "5958b185-8150-4aae-b53f-0c44771ddec5")
// 		context := appengine.NewContext(r)
// 		err := memcache.Delete(context, key)
// 		if err == nil {
// 			result[key] = true
// 			contentType, _ := url.QueryUnescape(parts[0])
// 			if imageTypes.MatchString(contentType) {
// 				thumbnailKey := key + thumbSuffix + filepath.Ext(parts[2])
// 				err := memcache.Delete(context, thumbnailKey)
// 				if err == nil {
// 					result[thumbnailKey] = true
// 				}
// 			}
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		b, err := json.Marshal(result)
// 		check(err)
// 		fmt.Fprintln(w, string(b))
// 	} else {
// 		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// ValidateType - Check for valid file types.
func (fi *FileInfo) ValidateType() (valid bool) {
	if acceptFileTypes.MatchString(fi.Type) {
		return true
	}
	fi.Error = "Filetype not allowed"
	return false
}

// ValidateSize - Check that file is within the enabled size range
func (fi *FileInfo) ValidateSize() (valid bool) {
	if fi.Size < minFileSize {
		fi.Error = "File is too small"
	} else if fi.Size > maxFileSize {
		fi.Error = "File is too big"
	} else {
		return true
	}
	return false
}

// func (fi *FileInfo) CreateUrls(r *http.Request, c context.Context) {
// 	u := &url.URL{
// 		Scheme: r.URL.Scheme,
// 		//Host:   appengine.DefaultVersionHostname(c),
// 		Host: "localhost:8080",
// 		Path: "/",
// 	}
// 	uString := u.String()
// 	fi.Url = uString + fi.Key
// 	fi.DeleteUrl = fi.Url
// 	fi.DeleteType = "DELETE"
// 	if fi.ThumbnailKey != "" {
// 		fi.ThumbnailUrl = uString + fi.ThumbnailKey
// 	}
// }

// SetKey - Creates and Sets the FileInfo Key property.
func (fi *FileInfo) SetKey(checksum uint32) {
	// fi.Key = escape(string(fi.Type)) + "/" +
	// 	escape(fmt.Sprint(checksum)) + "/" +
	// 	escape(string(fi.Name))
	fi.Key = path.Join(escape(string(fi.Type)), escape(fmt.Sprint(checksum)), escape(string(fi.Name)))
}

// func (fi *FileInfo) createThumb(buffer *bytes.Buffer, c context.Context) {
// 	if imageTypes.MatchString(fi.Type) {
// 		src, _, err := image.Decode(bytes.NewReader(buffer.Bytes()))
// 		check(err)
// 		filter := gift.New(gift.ResizeToFit(
// 			thumbMaxWidth,
// 			thumbMaxHeight,
// 			gift.LanczosResampling,
// 		))
// 		dst := image.NewNRGBA(filter.Bounds(src.Bounds()))
// 		filter.Draw(dst, src)
// 		buffer.Reset()
// 		bWriter := bufio.NewWriter(buffer)
// 		switch fi.Type {
// 		case "image/jpeg", "image/pjpeg":
// 			err = jpeg.Encode(bWriter, dst, nil)
// 		case "image/gif":
// 			err = gif.Encode(bWriter, dst, nil)
// 		default:
// 			err = png.Encode(bWriter, dst)
// 		}
// 		check(err)
// 		bWriter.Flush()
// 		thumbnailKey := fi.Key + thumbSuffix + filepath.Ext(fi.Name)
// 		item := &memcache.Item{
// 			Key:   thumbnailKey,
// 			Value: buffer.Bytes(),
// 		}
// 		err = memcache.Set(c, item)
// 		check(err)
// 		fi.ThumbnailKey = thumbnailKey
// 	}
// }
