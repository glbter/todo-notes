

package postgres

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"todoNote/internal/model"
)

const dbUrl = "postgres://postgres:123@localhost:5432/todoNote"

//func usersInit(t *testing.T) {
//	c := connect(t)
//	c.Exec(context.Background(), `insert into users
//    (id, name, time_zone, password_hash)
//values
//    (1, 'user1', 'UTC', 'user'),
//    (2, 'user2', 'UTC', 'user'),
//    (3, 'user3', 'UTC', 'user'),
//    (4, 'user4', 'UTC', 'user'),
//    (5, 'user5', 'UTC', 'user'),
//    (6, 'user6', 'UTC', 'user'),
//    (7, 'user7', 'UTC', 'user'),
//    (8, 'user8', 'UTC', 'user'),
//    (9, 'user9', 'UTC', 'user');`)
//	c.Close(context.Background())
//}

//func usersClean(t *testing.T) {
//	c := connect(t)
//	c.Exec(context.Background(), `truncate table users cascade;`)
//	c.Close(context.Background())
//}


func TestRepoUser_Insert(t *testing.T) {
	tts := []struct{
		caseName string
		hasError bool
	} {
		{
			caseName: "success",
			hasError: false,
		},
		{
			caseName: "failure",
			hasError: true,
		},
	}
	for _, tt := range tts {
		t.Run(tt.caseName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoUser(conn)

			u := model.NewUser(0, "insertUser", []byte("user"), model.UTC)
			id, err := r.Insert(context.Background(), u)
			assert.Equal(t, tt.hasError, err != nil)
			assert.NotEqual(t, 0, id)

			if !tt.hasError {
				var cid model.Id
				err = conn.QueryRow(context.Background(),
					`select id from users where id = $1`,
					id).Scan(&cid)
				assert.Nil(t, err)
			}
		})
	}

	//t.Run("clean up", func(t *testing.T) {
	//	usersClean(t)
	//})
}

func TestRepoUser_GetById(t *testing.T) {
	//t.Run("init", func(t *testing.T) {
	//	usersInit(t)
	//})

	tts := []struct{
		testName string
		hasErr bool
	}{
		{testName: "success", hasErr: false},
		{testName: "fail", hasErr: true},
	}

	for _, tt := range tts {
		t.Run(tt.testName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoUser(conn)

			//for id := 1; id < 10; id++ {
			id := 7
			user, err := r.GetById(context.Background(), model.Id(id))
			assert.Equal(t, tt.hasErr, err != nil)
			if err == nil {
				assert.Equal(t, model.Id(id), user.Id)
				assert.Equal(t, fmt.Sprintf("user%v", id), user.Name)
			}
			//}
		})

		//usersClean(t)
	}
}

func TestRepoUser_GetByUserName(t *testing.T) {
	//t.Run("init", func(t *testing.T) {
	//	usersInit(t)
	//})

	tts := []struct{
		testName string
		hasErr bool
	} {
		{testName: "success", hasErr: false},
		{testName: "fail", hasErr: true},
	}
	for _, tt := range tts {
		t.Run(tt.testName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoUser(conn)

			//for id := 1; id < 10; id++ {
			id := 8
			user, err := r.GetByUserName(context.Background(), fmt.Sprintf("user%v", id))
			assert.Equal(t, tt.hasErr, err != nil)
			if err == nil {
				assert.Equal(t, model.Id(id), user.Id)
			}
			//}

			//usersClean(t)
		})
	}
}

func TestRepoUser_Update(t *testing.T) {
	//t.Run("init", func(t *testing.T) {
	//	usersInit(t)
	//})

	t.Run("success", func(t *testing.T) {
		conn := connect(t)
		defer conn.Close(context.Background())

		r := NewRepoUser(conn)

		//for id := 1; id < 10; id++ {
		id := 6
			var user model.User
			err := conn.QueryRow(context.Background(),
				`select id, name, time_zone, password_hash from users where id = $1`,
				id).Scan(&user.Id, &user.Name, &user.TimeZone, &user.PasswordHash)
			assert.Nil(t, err)
			user.TimeZone = model.UTCp3
			err = r.Update(context.Background(), &user)
			assert.Nil(t, err)

			var updUser model.User
			err = conn.QueryRow(context.Background(),
				`select id, name, time_zone, password_hash from users where id = $1`,
				id).Scan(&updUser.Id, &updUser.Name, &updUser.TimeZone, &updUser.PasswordHash)
			assert.Nil(t, err)
			assert.Equal(t, user, updUser)
		//}
	})

	//t.Run("clean up", func(t *testing.T) {
	//	usersClean(t)
	//})
}

func TestRepoUser_Delete(t *testing.T) {
	//t.Run("set up", func(t *testing.T) {
	//	usersInit(t)
	//})

	tts := []struct{
		testName string
		hasErr bool
	} {
		{testName: "success", hasErr: false},
		{testName: "fail", hasErr: true},
	}
	for _, tt := range tts {
		t.Run(tt.testName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoUser(conn)
			i := 9
			//for i := 1; i < 10; i++ {
				err := r.Delete(context.Background(), model.Id(i))
				assert.Equal(t, tt.hasErr, err != nil)
			//}
		})
	}
	//t.Run("clean up", func(t *testing.T) {
	//	usersClean(t)
	//})
}


func TestRepoUser_Insert_GetById_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		conn := connect(t)
		defer conn.Close(context.Background())

		r := NewRepoUser(conn)

		u := model.NewUser(0, "user", []byte("user"), model.UTC)
		id, err := r.Insert(context.Background(), u)
		assert.Nil(t, err)
		assert.NotEqual(t, 0, id)

		user, err := r.GetById(context.Background(), id)
		assert.Nil(t, err)
		u.Id = id
		assert.Equal(t, u, user)

		err = r.Delete(context.Background(), id)
		assert.Nil(t, err)

		_, err = r.GetById(context.Background(), id)
		assert.NotNil(t, err)
	})
}

func TestRepoUser_Insert_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		conn := connect(t)
		defer conn.Close(context.Background())

		r := NewRepoUser(conn)

		for i := 0; i < 7; i++ {
			u := model.NewUser(0, fmt.Sprintf("user%v",i), []byte("user"), model.UTC)
			id, err := r.Insert(context.Background(), u)

			assert.Nil(t, err)
			assert.NotEqual(t, 0, id)

			defer func(id model.Id) {
				err = r.Delete(context.Background(), id)
				assert.Nil(t, err)
			} (id)
		}
	})
}

func TestRepoUser_GetByName_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		conn := connect(t)
		defer conn.Close(context.Background())

		r := NewRepoUser(conn)

		u := model.NewUser(0, "user", []byte("user"), model.UTC)
		id, err := r.Insert(context.Background(), u)
		assert.Nil(t, err)
		assert.NotEqual(t, 0, id)

		user, err := r.GetByUserName(context.Background(), u.Name)
		assert.Nil(t, err)
		u.Id = id
		assert.Equal(t, u, user)

		user.TimeZone = model.UTCp3

		err = r.Update(context.Background(), user)
		assert.Nil(t, err)

		updUsr, _ := r.GetById(context.Background(), id)
		assert.NotEqual(t, u, updUsr)
		assert.Equal(t, user, updUsr)

		r.Delete(context.Background(), id)
	})
}
