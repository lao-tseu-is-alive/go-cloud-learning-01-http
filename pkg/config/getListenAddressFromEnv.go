package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type ErrorConfig struct {
	err error
	msg string
}

func (e *ErrorConfig) Error() string {
	return fmt.Sprintf("%s : %v", e.msg, e.err)
}

//GetListenAddrFromEnv returns a valid TCP/IP listening address string based on the values of environment variables :
//	SERVERIP : string containing a valid Ip Address to use for listening (defaultIP will be used if env is not defined)
//	PORT : int value between 1 and 65535 (defaultPort will be used if env is not defined)
func GetListenAddrFromEnv(defaultIP string, defaultPort int) (string, error) {
	srvIP := defaultIP
	srvPort := defaultPort

	var err error
	val, exist := os.LookupEnv("PORT")
	if exist {
		srvPort, err = strconv.Atoi(val)
		if err != nil {
			return "", &ErrorConfig{
				err: err,
				msg: "ERROR: CONFIG ENV PORT should contain a valid integer.",
			}
		}
		if srvPort < 1 || srvPort > 65535 {
			return "", &ErrorConfig{
				err: err,
				msg: "ERROR: CONFIG ENV PORT should contain an integer between 1 and 65535",
			}
		}
	}
	val, exist = os.LookupEnv("SERVERIP")
	if exist {
		srvIP = val
		if net.ParseIP(srvIP) == nil {
			return "", &ErrorConfig{
				err: err,
				msg: "ERROR: CONFIG ENV SERVERIP should contain a valid IP Address.",
			}
		}
	}
	return fmt.Sprintf("%s:%d", srvIP, srvPort), nil
}
