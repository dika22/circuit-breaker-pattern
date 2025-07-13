# Circuit Breaker Pattern

#### how to Run 
``` go run main.go ```

## Documentation
Jika API third-party error 3x, circuit breaker akan terbuka dan semua request berikutnya langsung gagal.
Setelah timeout (5 detik), circuit breaker akan kembali ke half-open dan mencoba lagi.