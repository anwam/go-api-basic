package client

import (
	"context"
	"fmt"

	sub "bitbucket.org/truedmp/dmpss-pub-subscription-client"
	"bitbucket.org/truedmp/dmpss-pub-subscription-client/operation"
)

func GetProducts(ctx context.Context) {
	c := sub.Client{
		APIKey: "sdfasfa",
	}

	req := &operation.GetProductsV3Request{
		ProductGroupCode: []string{"TVSNOW"},
		ProductCode:      []string{"TVSNOWGOLD"},
		Lang:             "th",
		Country:          "th",
		Page:             1,
		Limit:            10,
	}
	res, err := c.GetInternalProductsV3(ctx, req)
	if err != nil {
		fmt.Printf("error naja : %s", err.Error())
	}
	fmt.Printf("response naja : %v", res)
}
