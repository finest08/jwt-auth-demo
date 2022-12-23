package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/finest08/jwt-auth-demo/authentication"
	"github.com/finest08/jwt-auth-demo/model"
	"github.com/finest08/jwt-auth-demo/store"
	"github.com/finest08/jwt-auth-demo/tjwt"
)

type Person struct {
	Store *store.Store
}

func (p *Person) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	defer r.Body.Close()
	reqByt, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	}

	var per model.Authentication
	json.Unmarshal(reqByt, &per)

	pass, _ := p.Store.VerifyPerson(per.Email)

	match, _ := authentication.VPass(per.Password, pass.Password)

	if !match {
		w.WriteHeader(http.StatusForbidden)
		return
	} else {
		var r = "user"
		tkn, _ := tjwt.GenJWT(per.Email, r)

		rsp, err := json.Marshal(model.Token{Status: true, Token: tkn})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %v", err)))
		}
		w.Write(rsp)
	}
}

func (p *Person) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		tk := r.Header.Get("Authorization")
		token, _ := tjwt.VerifyJWT(tk)

		claims := token.Claims.(jwt.MapClaims)

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			email := claims["email"].(string)

			psn, err := p.Store.PersonDetail(email)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %v", err)))
			}
			rsp, _ := json.Marshal(model.AuthResponse{Message: "Valid token", Status: true, ID: psn.ID})
			w.Write(rsp)
		}
		return
	}
}

func (p *Person) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	defer r.Body.Close()
	reqByt, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	}
	var per model.PersonCreate
	json.Unmarshal(reqByt, &per)

	per.ID = uuid.New().String()
	per.Date = time.Now()

	per.Password, _ = authentication.GenHashPass(per.Password)
	p.Store.AddPerson(per)
}

func (p *Person) Query(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		tk := r.Header.Get("Authorization")
		token, _ := tjwt.VerifyJWT(tk)

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			fn := r.URL.Query().Get("fn")
			ln := r.URL.Query().Get("ln")
			st := r.URL.Query().Get("st")
			lmtStr := r.URL.Query().Get("lmt")
			skipStr := r.URL.Query().Get("off")
			lmt, _ := strconv.ParseInt(lmtStr, 10, 64)
			skip, _ := strconv.ParseInt(skipStr, 10, 64)

			ppl, err := p.Store.GetPeople(fn, ln, st, &lmt, &skip)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %v", err)))
			}

			rspByt, err := json.Marshal(ppl)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %v", err)))
			}
			w.Write(rspByt)
		}
	}
}

func (p *Person) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	id := chi.URLParam(r, "id")

	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		tk := r.Header.Get("Authorization")
		token, _ := tjwt.VerifyJWT(tk)

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			psn, err := p.Store.GetPerson(id)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %v", err)))
			}
			data := model.Person{ID: psn.ID, GivenName: psn.GivenName, FamilyName: psn.FamilyName, Email: psn.Email, Phone: psn.Phone}
			rspByt, err := json.Marshal(data)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %v", err)))
			}
			w.Write(rspByt)
		}
		return
	}
}

func (p *Person) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	id := chi.URLParam(r, "id")

	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		tk := r.Header.Get("Authorization")
		token, _ := tjwt.VerifyJWT(tk)

		claims := token.Claims.(jwt.MapClaims)

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			defer r.Body.Close()
			reqByt, err := io.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("err %v", err)))
			}
			email := claims["email"].(string)

			var psn model.PersonCreate
			psn, err = p.Store.GetPerson(id)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %v", err)))
			}

			if psn.Email != email {
				w.Write([]byte("Unauthorized"))
				return
			} else {
				json.Unmarshal(reqByt, &psn)
				p.Store.UpdatePerson(id, psn)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("error %v", err)))
				}
				w.Write([]byte(""))
			}
		}
		return
	}
}

func (p *Person) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	id := chi.URLParam(r, "id")

	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		tk := r.Header.Get("Authorization")
		token, _ := tjwt.VerifyJWT(tk)

		claims := token.Claims.(jwt.MapClaims)

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			email := claims["email"].(string)

			psn, err := p.Store.GetPerson(id)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %v", err)))
			}

			if email != psn.Email {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				p.Store.DeletePerson(psn.ID)
				w.Write([]byte("{}"))
			}
		}
	}
}
