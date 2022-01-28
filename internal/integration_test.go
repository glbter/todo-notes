//+build integration

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/server/http/dto"
)

const (
	host = "http://localhost:8080"
	users = "api/v1/users"
	login = "api/v1/login"
	notes = "api/v1/notes"
	content = "application/json"
	auth = "Authorization"
)

var (
	usersUrl = fmt.Sprintf("%v/%v", host, users)
	loginUrl = fmt.Sprintf("%v/%v", host, login)
	notesUrl = fmt.Sprintf("%v/%v", host, notes)
	token string
	noteId int64
)

func Test_User_Register(t *testing.T) {
	b := dto.UserRegistration{
		UserName: "new user",
		Password: "123",
		TimeZone: model.UTC,
	}

	js, _ := json.Marshal(b)
	t.Run("first attempt", func(t *testing.T) {
		resp, err := http.Post(usersUrl, content, bytes.NewReader(js))
		assert.Nil(t, err)


		var id dto.IdObject
		err = json.NewDecoder(resp.Body).Decode(&id)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, 0, id.Id)
		resp.Body.Close()
	})

	t.Run("with same username", func(t *testing.T) {
		resp, err := http.Post(usersUrl, content, bytes.NewReader(js))
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}


func Test_User_Login(t *testing.T) {
	t.Run("wrong user", func(t *testing.T) {
		b := dto.UserLogin{
			UserName: "wrong user",
			Password: "123",
		}

		js, _ := json.Marshal(b)
		resp, err := http.Post(loginUrl, content, bytes.NewReader(js))
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("successful user", func(t *testing.T) {
		b := dto.UserLogin{
			UserName: "new user",
			Password: "123",
		}

		js, _ := json.Marshal(b)
		resp, err := http.Post(loginUrl, content, bytes.NewReader(js))
		assert.Nil(t, err)
		defer resp.Body.Close()

		var tk dto.Token
		err = json.NewDecoder(resp.Body).Decode(&tk)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		token = tk.Type + " " + tk.Token
	})
}

func Test_User_Patch(t *testing.T) {
	t.Run("successful user", func(t *testing.T) {
		b := dto.UserUpdate{
			TimeZone: model.UTCp3,
		}

		js, _ := json.Marshal(b)
		req, err := http.NewRequest(http.MethodPatch, usersUrl, bytes.NewReader(js))
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		//resp, err := http.Post(loginUrl, content, bytes.NewReader(js))
		resp, err := (&http.Client{}).Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		fmt.Println(resp.Body)
	})
}

func Test_Note_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		b := dto.NewNote{
			Title: "new note",
			Text: "New note",
			Date: time.Now(),
		}

		js, _ := json.Marshal(b)

		req, err := http.NewRequest(http.MethodPost, notesUrl, bytes.NewReader(js))
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		resp, err := (&http.Client{}).Do(req)
		assert.Nil(t, err)


		var id dto.IdObject
		err = json.NewDecoder(resp.Body).Decode(&id)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEqual(t, 0, id.Id)
		resp.Body.Close()

		noteId = id.Id
	})
}


func Test_Note_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, notesUrl, nil)
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		resp, err := (&http.Client{}).Do(req)
		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})
}



func Test_Note_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v", notesUrl, noteId), nil)
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		resp, err := (&http.Client{}).Do(req)
		assert.Nil(t, err)

		var b model.Note
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, noteId, b.Id)
		resp.Body.Close()
	})

	t.Run("fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v", notesUrl, 200000000), nil)
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		resp, err := (&http.Client{}).Do(req)
		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	})
}

func Test_Note_Patch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		b := dto.NoteUpdate{
			Text: "more new text",
			Title: "another title",
		}

		js, _ := json.Marshal(b)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%v/%v", notesUrl, noteId), bytes.NewReader(js))
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		resp, err := (&http.Client{}).Do(req)
		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()
	})
}

func Test_Note_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/%v", notesUrl, noteId), nil)
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		resp, err := (&http.Client{}).Do(req)
		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()
	})
}


func Test_User_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, usersUrl, nil)
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, token)
		resp, err := (&http.Client{}).Do(req)


		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()
		fmt.Println(resp.Body)
	})

	t.Run("bad token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, usersUrl, nil)
		req.Header.Set("Content-Type", content)
		req.Header.Set(auth, "bad token")
		resp, err := (&http.Client{}).Do(req)


		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		resp.Body.Close()
	})
}

//q := req.URL.Query()
//q.Add(timezoneQueryParam, model.UTC)
//req.URL.RawQuery = q.Encode()

