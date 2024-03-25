package api

import (
	"encoding/gob"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gabrielluizsf/go-web3/core"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
)

type TransactionResponse struct {
	TransactionCount uint
	Hashes           []string
}

type APIError struct {
	Error string
}

type Block struct {
	Hash          string
	Version       uint32
	DataHash      string
	PrevBlockHash string
	Height        uint32
	Timestamp     int64
	Validator     string
	Signature     string

	TransactionResponse
}

type ServerConfig struct {
	Logger     log.Logger
	ListenAddr string
}

type Server struct {
	TransactionChan chan *core.Transaction
	ServerConfig
	bc *core.Blockchain
}

func NewServer(cfg ServerConfig, bc *core.Blockchain, TransactionChan chan *core.Transaction) *Server {
	return &Server{
		ServerConfig:    cfg,
		bc:              bc,
		TransactionChan: TransactionChan,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	e.GET("/block/:hashorid", s.handleGetBlock)
	e.GET("/Transaction/:hash", s.handleGetTransaction)
	e.POST("/Transaction", s.handlePostTransaction)

	return e.Start(s.ListenAddr)
}

func (s *Server) handlePostTransaction(c echo.Context) error {
	transaction := &core.Transaction{}
	if err := gob.NewDecoder(c.Request().Body).Decode(transaction); err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}
	s.TransactionChan <- transaction

	return nil
}

func (s *Server) handleGetTransaction(c echo.Context) error {
	hash := c.Param("hash")

	b, err := hex.DecodeString(hash)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	transaction, err := s.bc.GetTransactionByHash(types.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, transaction)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	hashOrID := c.Param("hashorid")

	height, err := strconv.Atoi(hashOrID)
	// If the error is nil we can assume the height of the block is given.
	if err == nil {
		block, err := s.bc.GetBlock(uint32(height))
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, intoJSONBlock(block))
	}

	// otherwise assume its the hash
	b, err := hex.DecodeString(hashOrID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	block, err := s.bc.GetBlockByHash(types.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, intoJSONBlock(block))
}

func intoJSONBlock(block *core.Block) Block {
	transactionResponse := TransactionResponse{
		TransactionCount: uint(len(block.Transactions)),
		Hashes:           make([]string, len(block.Transactions)),
	}

	for i := 0; i < int(transactionResponse.TransactionCount); i++ {
		transactionResponse.Hashes[i] = block.Transactions[i].Hash(core.TransactionHasher{}).String()
	}

	return Block{
		Hash:                block.Hash(core.BlockHasher{}).String(),
		Version:             block.Header.Version,
		Height:              block.Header.Height,
		DataHash:            block.Header.DataHash.String(),
		PrevBlockHash:       block.Header.PrevBlockHash.String(),
		Timestamp:           block.Header.Timestamp,
		Validator:           block.Validator.Address().String(),
		Signature:           block.Signature.String(),
		TransactionResponse: transactionResponse,
	}
}
