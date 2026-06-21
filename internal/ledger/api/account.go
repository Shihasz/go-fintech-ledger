package api

import (
	"database/sql"
	"net/http"

	"github.com/Shihasz/go-fintech-ledger/internal/common/token"
	"github.com/Shihasz/go-fintech-ledger/internal/ledger/db"
	"github.com/gin-gonic/gin"
)

// createAccountRequest defines the JSON body we expect from the user.
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	// Parse and validate the JSON body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the authenticated payload from middleware context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Map the request to the sqlc database parameters
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0, // Initial balance is always 0
	}

	// Call the database function
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the newly created account
	ctx.JSON(http.StatusOK, account)
}

// getAccountRequest defines the URI parameters we expect.
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	// Parse and validate the URI ID parameter
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch from the database
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the account
	ctx.JSON(http.StatusOK, account)
}
