// Package models provides ...
package models

type ErrorModel struct {
	Code     int         `json:"code"`
	Err      string      `json:"err"`
	Desc     string      `json:"desc"`
	Internal interface{} `json:"internal"`
}
