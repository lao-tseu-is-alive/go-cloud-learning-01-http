#!/usr/bin/python3
# urlEncode.py is a small script that allows to url encode the given first argument
# typically it is a good practice to url encode your password before trying a dsn db connections
import sys
import urllib.parse
if len(sys.argv) < 2:
    print("ERROR : I expect first argument to be the string to url encode !")
    sys.exit(1)
givenString=sys.argv[1]
print(urllib.parse.quote(givenString))
