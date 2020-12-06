package app

import (
	"context"
	"encoding/json"
	"go2hw31/cmd/service/app/dto"
	"go2hw31/cmd/service/app/middleware/authenticator"
	"go2hw31/cmd/service/app/middleware/authorizator"
	"go2hw31/cmd/service/app/middleware/identificator"
	"go2hw31/pkg/business"
	"go2hw31/pkg/security"
	"go2hw31/pkg/web"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// ------------------------------------------------------------
type Server struct {
	securitySvc *security.Service
	businessSvc *business.Service
	router      chi.Router
}

func NewServer(securitySvc *security.Service, businessSvc *business.Service, router chi.Router) *Server {
	return &Server{
		securitySvc: securitySvc,
		businessSvc: businessSvc,
		router:      router,
	}
}

func (s *Server) Init() error {
	identificatorMd := identificator.Identificator
	authenticatorMd := authenticator.Authenticator(
		identificator.Identifier, s.securitySvc.UserDetails,
	)
	// функция-связка между middleware и security service
	// (для чистоты security service ничего не знает об http)
	roleChecker := func(ctx context.Context, roles ...string) bool {
		userDetails, err := authenticator.Authentication(ctx)
		if err != nil {
			return false
		}
		return s.securitySvc.HasAnyRole(ctx, userDetails, roles...)
	}
	//adminRoleMd := authorizator.Authorizator(roleChecker, security.RoleAdmin)
	//userRoleMd := authorizator.Authorizator(roleChecker, security.RoleUser)

	s.router.Get("/echo", s.handleEcho)
	s.router.Post("/api/users", s.handleRegister2)
	s.router.Post("/tokens", s.handleTokens)

	//s.router.With(identificatorMd, authenticatorMd, adminRoleMd).Get("/cardsAdmin", s.handleCardsAdmin)
	//s.router.With(identificatorMd, authenticatorMd, userRoleMd).Get("/cards", s.handleCardsUser)

	s.router.With(identificatorMd, authenticatorMd, authorizator.Authorizator(roleChecker, security.RoleAdmin, security.RoleUser)).Get("/cards", s.handleGetCards)

	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

// ------------------------------------------------------------
func (s *Server) handleEcho(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write([]byte("this is echo page"))
	if err != nil {
		log.Print(err)
	}
}

// ------------------------------------------------------------
// получить профиль юзера из контекста
func getProfile(ctx context.Context) (*security.UserDetails, error) {
	profile, err := authenticator.Authentication(ctx)
	if err != nil {
		return nil, err
	}
	profile2, ok := profile.(*security.UserDetails)
	if !ok {
		log.Print(err)
	}
	return profile2, nil
}

// ------------------------------------------------------------
type LoginPasswordIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) handleRegister2(writer http.ResponseWriter, request *http.Request) {
	var qparams LoginPasswordIn
	err := json.NewDecoder(request.Body).Decode(&qparams)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	login := qparams.Login
	if login == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	password := qparams.Password
	if password == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := s.securitySvc.Register2(request.Context(), login, password)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := map[string]int64{"id": *userID}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(respJSON)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleTokens(writer http.ResponseWriter, request *http.Request) {
	var qparams LoginPasswordIn
	err := json.NewDecoder(request.Body).Decode(&qparams)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	login := qparams.Login
	pass := qparams.Password
	token, err := s.securitySvc.Login(request.Context(), login, pass)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &dto.TokenDTO{Token: token}
	respBody, err := json.Marshal(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(respBody)
	if err != nil {
		log.Print(err)
	}

}

// ------------------------------------------------------------
/*
func (s *Server) handleCardsAdmin(writer http.ResponseWriter, request *http.Request) {
	cards, err := s.businessSvc.GetAllCards(request.Context())
	if err != nil {
		log.Print(err)
	}
	err = web.WriteAsJSON(writer, cards)
	if err != nil {
		log.Print(err)
	}

}

func (s *Server) handleCardsUser(writer http.ResponseWriter, request *http.Request) {
	profile, err := getProfile(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	userCards, err := s.businessSvc.GetUserCards(request.Context(), profile.ID)
	if err != nil {
		log.Print(err)
	}
	err = web.WriteAsJSON(writer, userCards)
	if err != nil {
		log.Print(err)
	}
}
*/

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// общий хендлер
func (s *Server) handleGetCards(writer http.ResponseWriter, request *http.Request) {
	profile, err := getProfile(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// admin
	if contains(profile.Roles, security.RoleAdmin) {
		cards, err := s.businessSvc.GetAllCards(request.Context())
		if err != nil {
			log.Print(err)
		}
		err = web.WriteAsJSON(writer, cards)
		if err != nil {
			log.Print(err)
		}
		return
	}

	// user
	if contains(profile.Roles, security.RoleUser) {
		userCards, err := s.businessSvc.GetUserCards(request.Context(), profile.ID)
		if err != nil {
			log.Print(err)
		}
		err = web.WriteAsJSON(writer, userCards)
		if err != nil {
			log.Print(err)
		}
		return
	}

}
