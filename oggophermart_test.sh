#!/bin/zsh
go build -buildvcs=false -o oggophermart
cd accrual_bin && chmod +x accrual_linux_amd64
gophermarttest \
            -test.v -test.run=^TestGophermart$ \
            -gophermart-binary-path=./oggophermart \
            -gophermart-host=localhost \
            -gophermart-port=8080 \
            -gophermart-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" \
            -accrual-binary-path=accrual_bin/accrual_linux_amd64 \
            -accrual-host=localhost \
            -accrual-port=8081 \
            -accrual-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"
