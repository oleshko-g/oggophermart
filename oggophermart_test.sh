#!/bin/zsh
go build -buildvcs=false -o oggophermart
(cd accrual_bin && chmod +x accrual_darwin_arm64) &&
gophermarttest \
-test.v -test.run=^TestGophermart$ \
-gophermart-binary-path=./oggophermart \
-gophermart-host=localhost \
-gophermart-port=8080 \
-gophermart-database-uri="postgres://gennadyoleshko:@localhost:5432/oggophermart?sslmode=disable" \
-accrual-binary-path=accrual_bin/accrual_darwin_arm64 \
-accrual-host=localhost \
-accrual-port=8081 \
-accrual-database-uri="postgres://gennadyoleshko@localhost:5432/praktikum?sslmode=disable"
