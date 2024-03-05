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

## CONFIG TS ADAPTER
1. Pastikan config way4 manager telah terconfig.
2. Copy `/way4ts/appserver/applications/TS1A/webapps/TS1A/WEB-INF/inv/scripts/FrAlfamartAdapter.s.xml` ke platform TS yang akan dideploy.
3. Copy `way4ts/appserver/applications/TS1A/webapps/TS1A/WEB-INF/conf/application/FrameAlfamart.groovy` ke platform TS yang akan dideploy.
4. Jalankan service `FrAlfamartAdapter` melalui m2_web_console TS


## HOW TO START DEV ENVIRONMENT ##

1. Pastikan ts adapter dan config way4 manager telah terdeploy (lihat bagian DEPLOY TS ADAPTER)
2. Pastikan DB alfamart_channel telah terdeploy
3. Pointing `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_ADDRESS`, `DB_PORT` ke db alfamart_channel
4. Untuk menjalankan program ini pastikan value pada CI/CD variable telah terconfig dengan lengkap.
5. Set Environments pada CI/CD ke `staging` dan flags protected diset `true`
6. Jalankan stages pipeline `build` (hanya jika ada perubahan pada codes)
7. Jalankan stages pipeline `deploy_staging`
8. Jika semua stages sukses maka lakukan test request dengan menggunakan `LOADBALANCER_IP` yang telah diset di di CI/CD variable.



## HOW TO START PROD ENVIRONMENT ##

1. Pastikan ts adapter dan config way4 manager telah terdeploy (lihat bagian DEPLOY TS ADAPTER)
2. Pastikan DB alfamart_channel telah terdeploy
3. Pointing `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_ADDRESS`, `DB_PORT` ke db alfamart_channel
4. Untuk menjalankan program ini pastikan value pada CI/CD variable telah terconfig dengan lengkap.
5. Set Environments pada CI/CD ke `production` dan flags protected diset `true`
6. Jalankan stages pipeline `build` (hanya jika ada perubahan pada codes)
7. Jalankan stages pipeline `deploy_production`
8. Jika semua stages sukses maka lakukan test request dengan menggunakan `LOADBALANCER_IP` yang telah diset di di CI/CD variable.