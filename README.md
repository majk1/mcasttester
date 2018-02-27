# UDP Multicast tester CLI application [![Build Status](https://travis-ci.org/majk1/mcasttester.svg?branch=master)](https://travis-ci.org/majk1/mcasttester)

## Install

```
go get github.com/majk1/mcasttester
```

## Usage

```
mcasttester <command> [options]
```

 * commands and options:  
   ```
   commands:
   
       send    - Use command "send" to send data to the specified multicast address
       
           options:
               -addr string    - Multicast address and port to send data (default "224.0.0.1:9999")
               -data string    - The string to send (default "DATA")
               -loop int       - Loop count (default is forever)
   
       receive - Use command "receive" to to listen on the specified multicast address
   
           options:
               -addr string    - Multicast address and port to listen on (default "224.0.0.1:9999")
               -loop int       - Loop count (default is forever)
   
   ```

