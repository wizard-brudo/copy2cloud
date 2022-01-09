package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	ERROR_ARGUMENTS            = errors.New(setTextColor("Error:no launch arguments", RED))
	ERROR_LOGIN                = errors.New(setTextColor("Error:wrong username and/or password", RED))
	ERROR_UNKNOWN_ARGUMENT     = errors.New(setTextColor("Error:unknown argument", RED))
	ERROR_JSON                 = errors.New(setTextColor("Error:not correct json", RED))
	ERROR_NOT_ENOUGH_ARGUMENTS = errors.New(setTextColor("Error:not enough arguments", RED))
	ERROR_UNKOWN_STYLE         = errors.New(setTextColor("Error:unknown style", RED))
	ERROR_WRITE_BIT_NOT_SET    = errors.New(setTextColor("Error:write permission bit is not set on this file ", RED))
	ERROR_STAT                 = errors.New(setTextColor("Error:unable to get stat", RED))
	ERROR_NO_PERMISSION        = errors.New(setTextColor("Error:you don't have permission to write to this directory", RED))
	ERROR_NOT_DIRECTORY        = errors.New(setTextColor("Error:path isn't a directory", RED))
	ERROR_PATH_NOT_EXISTS      = errors.New(setTextColor("Error:path doesn't exist", RED))
	ERROR_NO_TOKEN             = errors.New(setTextColor("Error:no token", RED))
	ERROR_TOO_MANY_ARGUMENTS   = errors.New(setTextColor("Error:too many arguments", RED))
	ERROR_NO_RESOURCES         = errors.New(setTextColor("Error:no resources", RED))
)

func CheckResponse(response http.Response) error {
	errApi := ErrorApi{}
	if response.StatusCode >= 400 {
		body_byte, _ := ioutil.ReadAll(response.Body)
		err := json.Unmarshal(body_byte, &errApi)
		if err != nil {
			fmt.Println(ERROR_JSON)
			os.Exit(1)
		}
		return errors.New(setTextColor(errApi.Error+":"+errApi.Description, RED))
	}
	return nil
}
