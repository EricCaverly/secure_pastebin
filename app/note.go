/*********************************
 *  File     : note.go
 *  Purpose  : Defines what a note is
 *  Authors  : Eric Caverly
 */

package main

type Note struct {
	Content        string `json:"content" redis:"content"`
	AllowedIPRange string `json:"allowed_ips" redis:"allowed_ips"`
}
