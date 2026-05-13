package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: apikeyctl [create|revoke]")
	}
	switch os.Args[1] {
	case "create":
		runCreate(os.Args[2:])
	case "revoke":
		runRevoke(os.Args[2:])
	default:
		log.Fatalf("unknown command %q", os.Args[1])
	}
}

func runCreate(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	baseURL := fs.String("base-url", "http://localhost:8080", "")
	adminToken := fs.String("admin-token", "", "")
	app := fs.String("app", "", "")
	env := fs.String("env", "", "")
	name := fs.String("name", "", "")
	createdBy := fs.String("created-by", "", "")
	var scopes multiFlag
	fs.Var(&scopes, "scope", "")
	fs.Parse(args)

	body, _ := json.Marshal(map[string]any{
		"app":        *app,
		"env":        *env,
		"name":       emptyStringNil(*name),
		"created_by": emptyStringNil(*createdBy),
		"scopes":     []string(scopes),
		"expires_at": nil,
	})
	req, err := http.NewRequest(http.MethodPost, *baseURL+"/admin/api-keys", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Token", *adminToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Fatalf("request failed: %s", resp.Status)
	}
	var out bytes.Buffer
	if _, err := out.ReadFrom(resp.Body); err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}

func runRevoke(args []string) {
	fs := flag.NewFlagSet("revoke", flag.ExitOnError)
	baseURL := fs.String("base-url", "http://localhost:8080", "")
	adminToken := fs.String("admin-token", "", "")
	id := fs.Int64("id", 0, "")
	fs.Parse(args)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/admin/api-keys/%d/revoke", *baseURL, *id), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Admin-Token", *adminToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Fatalf("request failed: %s", resp.Status)
	}
	var out bytes.Buffer
	if _, err := out.ReadFrom(resp.Body); err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}

type multiFlag []string

func (m *multiFlag) String() string { return "" }
func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func emptyStringNil(value string) any {
	if value == "" {
		return nil
	}
	return value
}
