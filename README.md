# fastTCPscan


>`./fastTCPscan -h` 
>Usage of ./fastTCPscan:
>- host string
>        Host o direccion IP a escanear (default "127.0.0.1")
>- range string
>    Rango de puertos a escanear: 22,80-443,1443,2000-3000,1-65535... (default "1-65535")
>- threads int
>        Número de hilos a usar (default 750)
>- timeout duration
>        Segundos por puerto (default 1s)


##example:

`go run fastTCPscan.go -host=10.10.10.28`

>`[*] Escaneando host 10.10.10.28 (Puertos: 1-65535)`
>
>`22: Abierto`
>`80: Abierto`

`./fastTCPscan -host=10.10.10.28 -range=1-1000`

>`[*] Escaneando host 10.10.10.28 (Puertos: 1-1000)`
>
>`22: Abierto`
>`80: Abierto`

