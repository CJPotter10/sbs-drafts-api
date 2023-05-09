package utils

import (
	"math/big"
	"strings"

	"github.com/CJPotter10/sbs-drafts-api/api"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DraftTokenContract struct {
	EthConn *ethclient.Client
	Api     *api.ApiCaller
}

var Contract *DraftTokenContract

func CreateEthConnection(contractAddress string, infuraEndpoint string) error {
	contractAddr := common.HexToAddress(contractAddress)

	conn, err := ethclient.Dial(infuraEndpoint)
	if err != nil {
		return err
	}

	contract, err := api.NewApiCaller(contractAddr, conn)
	if err != nil {
		return err
	}

	Contract.EthConn = conn
	Contract.Api = contract
	return nil
}

func (c *DraftTokenContract) GetOwnerOfToken(tokenId int) (string, error) {
	id := big.NewInt(int64(tokenId))
	owner, err := c.Api.OwnerOf(nil, id)
	if err != nil {
		return "", err
	}
	return strings.ToLower(owner.Hex()), nil
}

func (c *DraftTokenContract) GetNumTokensMinted() (int, error) {
	res, err := c.Api.NumTokensMinted(nil)
	if err != nil {
		return -1, nil
	}
	return int(res.Int64()), nil
}
