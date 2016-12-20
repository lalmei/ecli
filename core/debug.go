package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func DebugResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("An error occurred. Here is some debug info:\n\n")
		fmt.Println(">> RESPONSE HEADER")
		for k, v := range resp.Header {
			fmt.Printf("%s: %q\n", k, v[0])
		}
		fmt.Println()
		fmt.Println(">> RESPONSE BODY")
		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
		fmt.Println()
		fmt.Println(">> HTTP ERROR")
		return fmt.Errorf(resp.Status)
	}
	return nil
}
