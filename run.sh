#!/bin/bash


echo "⏫ Iniciando contenedor MySQL (mysql-flighthours)..."
sudo docker start mysql-flighthours


echo "⏳ Esperando a que el contenedor inicie..."
sleep 5


echo "🚀 Ejecutando la aplicación Go..."
go run /home/devban/Documents/Go/flighthours-api/cmd/main.go
