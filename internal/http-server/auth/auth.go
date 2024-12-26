package auth

// TODO: разобраться с авторизацией, первая попытка неудачна
/*
type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	JWT string `json:"jwt,omitempty"`
}

type LogServ interface {
	SelectUser(login string) (model.User, error)
}

type RegServ interface {
	SaveUser(login string, hashPass []byte) error
}

func Register(log *slog.Logger, reg RegServ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.auth.Register"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Response: resp.Error("failed to generate password hash"),
			})
			return
		}
		if err := reg.SaveUser(req.Login, passHash); err != nil {
			log.Info("failed to register user", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Response: resp.Error("failed to register user"),
			})
			return
		}
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})

	}

}

func Login(log *slog.Logger, login LogServ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.auth.Login"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		user, err := login.SelectUser(req.Login)
		if err != nil {
			log.Debug("failed to find user", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Response: resp.Error("invalid credentials"),
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword(user.HashPas, []byte(req.Password)); err != nil {
			log.Debug("invalid credentials", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Response: resp.Error("invalid credentials"),
			})
			return
		}

		token, err := jwt.NewToken(user, os.Getenv("SECRET"), 12*time.Hour)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Response: resp.Error("failed to generate token"),
			})
			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			JWT:      token,
		})
	}

}

*/
