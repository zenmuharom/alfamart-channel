## ALFAMART CHANNEL ##
---------------------------

## Structure Folder ##

```
.
├── ROOT 
├── api 
├── domain
├── dto
├── function
├── helper
├── logger
├── models
├── repo
├── service
├── src
├── tool
├── util
├── variable
├── app.env
├── example.env
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── start.sh
.
```


## HOW TO START DEV ENVIRONMENT ##

1. Pastikan DB alfamart_channel telah terdeploy
2. Pointing `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_ADDRESS`, `DB_PORT` ke db alfamart_channel
3. Untuk menjalankan program ini pastikan value pada CI/CD variable telah terconfig dengan lengkap.
4. Set Environments pada CI/CD ke `staging` dan flags protected diset `true`
5. Jalankan stages pipeline `build` (hanya jika ada perubahan pada codes)
6. Jalankan stages pipeline `deploy_staging`
7. Jika semua stages sukses maka lakukan test request dengan menggunakan `LOADBALANCER_IP` yang telah diset di di CI/CD variable.