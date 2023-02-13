package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/arthurshafikov/tz-faraway/lib/powtcp"
)

func MakeQuery(ctx context.Context) error {
	address := "localhost:3333"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address, nil)
	if err != nil {
		return err
	}

	connDialer, err := powtcp.NewConnDialer(powtcp.ConnDialerOptions{
		Address:              address,
		ReadTimeoutDuration:  time.Second * 5,
		WriteTimeoutDuration: time.Second * 5,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := connDialer.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()
	client := http.Client{Transport: &http.Transport{Dial: connDialer.Dial}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Printf("%s %s", res.Status, string(bodyRes))

	return nil
}
