package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/santinofajardo/simpleBank/db/sqlc"
	"github.com/santinofajardo/simpleBank/token"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) transfer(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	ok, fromAccount := server.validCurrencyAccount(ctx, req.FromAccountID, req.Currency)

	if !ok {
		return
	}

	if authPayload.UserName != fromAccount.Owner {
		err := errors.New("fromAccount mismatch with the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if ok, _ := server.validCurrencyAccount(ctx, req.ToAccountID, req.Currency); !ok {
		return
	}

	arg := db.TransferTransactionParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTransaction(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result, "error": nil})
}

func (server *Server) validCurrencyAccount(ctx *gin.Context, accountID int64, currency string) (bool, db.Account) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err, "data": nil})
			return false, db.Account{}
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": err, "data": nil})
		return false, db.Account{}
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err, "data": nil})
		return false, db.Account{}
	}
	return true, account
}
