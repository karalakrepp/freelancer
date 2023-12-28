package authService

// import (
// 	"context"
// 	"time"

// 	"github.com/karalakrepp/Golang/freelancer-project/database"
// 	errs "github.com/karalakrepp/Golang/freelancer-project/errors"
// 	"github.com/karalakrepp/Golang/freelancer-project/models"
// 	"github.com/karalakrepp/Golang/freelancer-project/token"
// 	"github.com/karalakrepp/Golang/freelancer-project/util"
// )

// type Manager interface {
// 	// unauthorized
// 	SignUp(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
// 	Login(ctx context.Context, req models.LoginUserRequest) (models.LoginUserResponse, error)

// 	// authorized

// 	DeleteAccount(ctx context.Context, req models.DeleteAccountRequest) (models.DeleteAccountResponse, error)
// }

// type authManager struct {
// 	config     util.Config
// 	tokenMaker token.Maker
// 	db         database.Storage
// }

// func NewManager(config util.Config, tokenMaker token.Maker, db database.Storage) Manager {
// 	return &authManager{
// 		config:     config,
// 		tokenMaker: tokenMaker,
// 		db:         db,
// 	}
// }

// func (a *authManager) SignUp(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
// 	hashedPassword, err := util.HashPassword(req.Password)
// 	if err != nil {
// 		return models.UserResponse{}, err
// 	}

// 	user, err := a.db.GetUserByUserName(ctx, req.Username)
// 	if err == nil {
// 		if time.Now().Before(user.DeletedAt.AddDate(0, 0, 15)) {
// 			return models.UserResponse{}, errs.ErrorAccountIsDeleted
// 		}
// 	}

// 	createUserParams := database.CreateUserParams{
// 		FirstName:      req.FirstName,
// 		Username:       req.Username,
// 		Email:          req.Email,
// 		HashedPassword: hashedPassword,
// 	}

// 	resp, err := a.db.CreateAccount(createUserParams)
// 	if err != nil {
// 		return models.UserResponse{}, err
// 	}

// 	return newUserResponse(resp.User), nil
// }

// // func (a *authManager) Login(ctx context.Context, req models.LoginUserRequest) (models.LoginUserResponse, error) {
// // 	params := database.GetUserParams{
// // 		ID:       req.UserID,
// // 		Username: req.Username,
// // 	}

// // 	user, err := a.db.GetUser(ctx, params)
// // 	if err != nil {
// // 		return models.LoginUserResponse{}, errs.ErrorNoUser
// // 	}

// // 	err = util.CheckPassword(req.Password, user.HashedPassword)
// // 	if err != nil {
// // 		return models.LoginUserResponse{}, errs.ErrorUnauthorized
// // 	}

// // 	if !user.IsActive {
// // 		return models.LoginUserResponse{}, errs.ErrorAccountIsDeactivated
// // 	}

// // 	accessToken, accessPayload, err := a.tokenMaker.CreateToken(
// // 		user.ID,
// // 		a.config.ACCESS_TOKEN_DURATION,
// // 	)
// // 	if err != nil {
// // 		return models.LoginUserResponse{}, err
// // 	}

// // 	refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(
// // 		user.ID,
// // 		a.config.REFRESH_TOKEN_DURATION,
// // 	)
// // 	if err != nil {
// // 		return models.LoginUserResponse{}, err
// // 	}

// // 	sessionParams := database.CreateSessionParams{
// // 		ID:           refreshPayload.ID,
// // 		UserID:       user.ID,
// // 		RefreshToken: refreshToken,
// // 		UserAgent:    ctx.Value(models.UserAgent).(string),
// // 		ClientIp:     ctx.Value(models.RemoteAddress).(string),
// // 		IsBlocked:    false,
// // 		ExpiresAt:    refreshPayload.ExpiredAt,
// // 	}
// // 	session, err := a.db.CreateSession(ctx, sessionParams)
// // 	if err != nil {
// // 		return models.LoginUserResponse{}, err
// // 	}

// // 	resp := models.LoginUserResponse{
// // 		SessionID:             session.ID,
// // 		AccessToken:           accessToken,
// // 		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
// // 		RefreshToken:          refreshToken,
// // 		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
// // 		User:                  newUserResponse(user),
// // 	}

// // 	return resp, nil
// // }

// // func (a *authManager) RenewAccess(ctx context.Context, req models.RenewAccessTokenRequest) (models.RenewAccessTokenResponse, error) {
// // 	refreshPayload, err := a.tokenMaker.VerifyToken(req.RefreshToken)
// // 	if err != nil {
// // 		return models.RenewAccessTokenResponse{}, err
// // 	}

// // 	session, err := a.db.GetSession(ctx, refreshPayload.ID)
// // 	if err != nil {
// // 		return models.RenewAccessTokenResponse{}, err
// // 	}

// // 	if session.IsBlocked {
// // 		return models.RenewAccessTokenResponse{}, errs.ErrorUnauthorized
// // 	}

// // 	if session.UserID != refreshPayload.UserID {
// // 		return models.RenewAccessTokenResponse{}, errs.ErrorUnauthorized
// // 	}

// // 	if session.RefreshToken != req.RefreshToken {
// // 		return models.RenewAccessTokenResponse{}, errs.ErrorUnauthorized
// // 	}

// // 	if time.Now().After(session.ExpiresAt) {
// // 		return models.RenewAccessTokenResponse{}, errs.ErrorExpiredSession
// // 	}

// // 	accessToken, accessPayload, err := a.tokenMaker.CreateToken(
// // 		refreshPayload.UserID,
// // 		a.config.ACCESS_TOKEN_DURATION,
// // 	)
// // 	if err != nil {
// // 		return models.RenewAccessTokenResponse{}, err
// // 	}

// // 	resp := models.RenewAccessTokenResponse{
// // 		UserID:               refreshPayload.UserID,
// // 		AccessToken:          accessToken,
// // 		AccessTokenExpiresAt: accessPayload.ExpiredAt,
// // 	}

// // 	return resp, nil
// // }

// // func (a *authManager) VerifyEmail(ctx context.Context, id int64, secret_code string) (models.VerifyEmailResponse, error) {
// // 	txResult, err := a.db.VerifyEmailTx(ctx, database.VerifyEmailTxParams{
// // 		ID:         id,
// // 		SecretCode: secret_code,
// // 	})
// // 	if err != nil {
// // 		return models.VerifyEmailResponse{}, err
// // 	}

// // 	resp := models.VerifyEmailResponse{
// // 		UserID:     txResult.VerifyEmail.UserID,
// // 		IsVerified: txResult.User.IsEmailVerified,
// // 	}

// // 	return resp, nil
// // }

// // func (a *authManager) ResendEmail(ctx context.Context, req models.ResendEmailRequest) {
// // 	a.worker.EnqueueSendVerifyEmail(database.GetUserParams{
// // 		ID:       req.UserID,
// // 		Username: req.Username,
// // 	})
// // }

// // func (a *authManager) DeactivateAccount(ctx context.Context, req models.DeactivateAccountRequest) (models.DeactivateAccountResponse, error) {
// // 	authPayload := ctx.Value(models.AuthorizationPayload).(*token.Payload)
// // 	if authPayload.UserID != req.UserID {
// // 		return models.DeactivateAccountResponse{}, errs.ErrorNotAuthorized
// // 	}

// // 	user, err := a.db.GetUserById(ctx, authPayload.UserID)
// // 	if err != nil {
// // 		return models.DeactivateAccountResponse{}, err
// // 	}

// // 	err = util.CheckPassword(req.Password, user.HashedPassword)
// // 	if err != nil {
// // 		return models.DeactivateAccountResponse{}, errs.ErrorUnauthorized
// // 	}

// // 	_, err = a.db.UpdateUser(ctx, database.UpdateUserParams{
// // 		ID: authPayload.UserID,
// // 		IsActive: pgtype.Bool{
// // 			Bool:  false,
// // 			Valid: true,
// // 		},
// // 		DeactivatedAt: pgtype.Timestamptz{
// // 			Time:  time.Now(),
// // 			Valid: true,
// // 		},
// // 	})
// // 	if err != nil {
// // 		return models.DeactivateAccountResponse{}, err
// // 	}

// // 	return models.DeactivateAccountResponse{
// // 		Message: "Account Deactivated Successfully",
// // 	}, nil
// // }

// // // core deletion logic should move to a worker
// // func (a *authManager) DeleteAccount(ctx context.Context, req models.DeleteAccountRequest) (models.DeleteAccountResponse, error) {
// // 	authPayload := ctx.Value(models.AuthorizationPayload).(*token.Payload)
// // 	if authPayload.UserID != req.UserID {
// // 		return models.DeleteAccountResponse{}, errs.ErrorNotAuthorized
// // 	}

// // 	user, err := a.db.GetUserById(ctx, authPayload.UserID)
// // 	if err != nil {
// // 		return models.DeleteAccountResponse{}, err
// // 	}

// // 	err = util.CheckPassword(req.Password, user.HashedPassword)
// // 	if err != nil {
// // 		return models.DeleteAccountResponse{}, errs.ErrorUnauthorized
// // 	}

// // 	_, err = a.db.UpdateUser(ctx, database.UpdateUserParams{
// // 		ID: authPayload.UserID,
// // 		IsDeleted: pgtype.Bool{
// // 			Bool:  true,
// // 			Valid: true,
// // 		},
// // 	})
// // 	if err != nil {
// // 		return models.DeleteAccountResponse{}, err
// // 	}

// // 	a.worker.EnqueueDeleteOperation(authPayload.UserID)

// // 	return models.DeleteAccountResponse{
// // 		Message: "Account Deleted Successfully",
// // 	}, nil
// // }

// // // helper functions
// // func newUserResponse(user database.User) models.UserResponse {
// // 	return models.UserResponse{
// // 		ID:                user.ID,
// // 		Username:          user.Username,
// // 		Email:             user.Email,
// // 		PasswordChangedAt: user.PasswordChangedAt,
// // 		CreatedAt:         user.CreatedAt,
// // 	}
// // }
