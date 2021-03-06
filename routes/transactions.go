package routes

import (
	ethStreamerServices "github.com/AkhiraChain/eth-stalker/services"
	"github.com/gin-gonic/gin"
)

// SetTransactionsRoutes defines all transactions routes
func SetTransactionsRoutes(router *gin.Engine) {
	router.GET("/tx-start", ethStreamerServices.StreamBlockTransactions)
	router.GET("/transactions", ethStreamerServices.GetLatestBlockTransactions)
}
