# Circuit Breaker Pattern

#### how to Run 
``` go run main.go ```

## Documentation
Karena kita gunakan https://httpstat.us/503, hasilnya:

Retry 3x → tetap gagal

Circuit breaker akan "Open" setelah 3 kegagalan → langsung tolak permintaan berikutnya

Setelah 10 detik, "Half-Open", dan coba lagi