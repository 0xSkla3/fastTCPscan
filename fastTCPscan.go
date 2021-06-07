package main

import (
	"fmt"
	"flag"
	"context"
	"strings"
	"strconv"
	"log"
	"sync"
	"time"
	"net"
)

var (
	host = flag.String("host","127.0.0.1","Host o direccion IP a escanear")
	ports = flag.String("range","1-65535","Rango de puertos a escanear: 22,80-443,1443,2000-3000,1-65535...")
	threads = flag.Int("threads",750,"Número de hilos a usar")
	timeout = flag.Duration("timeout",1*time.Second,"Segundos por puerto")
)

func processRange(ctx context.Context, r string) chan int {
	c := make(chan int) // c <- elemento, 'c' es un conducto
	done := ctx.Done()

	go func(){
		defer close(c)
		blocks := strings.Split(r, ",")

		for _, block := range blocks{
			rg := strings.Split(block, "-")
			var minPort, maxPort int
			var err error

			minPort, err = strconv.Atoi(rg[0])
			if(err != nil){
				log.Print("No ha sido posible interpretar el rango: ", block)
				continue
			}

			if len(rg) == 1{
				maxPort = minPort
			} else{
				maxPort, err = strconv.Atoi(rg[1])
				if err != nil {
					log.Print("No ha sido posible interpretar el rango: ", block)
					continue
				}
			}
			for port := minPort; port <= maxPort; port++ {
				select {
				case c <- port:
				case <- done:
					return
				}
			}
		}
	}()
	return c

}

func scanPort(port int) string {
	addr := fmt.Sprintf("%s:%d", *host, port) //ip:puerto
	conn, err := net.DialTimeout("tcp", addr, *timeout)
	if err != nil{
		return fmt.Sprintf("%d: %s", port, err.Error())
	}

	conn.Close()

	return fmt.Sprintf("%d: Abierto", port)
}

func scanPorts(ctx context.Context, in <- chan int) chan string{
	out := make(chan string)
	done := ctx.Done()
	var wg sync.WaitGroup
	wg.Add(*threads)

	for i := 0; i < *threads; i++ {
		go func(){
			defer wg.Done()
			for {
				select{
					case port, ok := <- in:
						if !ok {
							return
						}
						s := scanPort(port)
						select {
						case out <- s:
						case <- done:
							return
						}
					case <- done:
						return
				}
			}
		}()
		go func(){
			wg.Wait()
			close(out)
		}()

	}
	return out
}


func main(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()
	fmt.Printf("\n [*] Escaneando host %s (Puertos: %s)\n\n", *host, *ports)
	pR := processRange(ctx, *ports)
	sP := scanPorts(ctx, pR)

	for port := range sP{
		if strings.HasSuffix(port, ": Abierto"){
			fmt.Println(port)
		}
	}

}
