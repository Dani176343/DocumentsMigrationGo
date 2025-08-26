// Arquivo: main.go
package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	// 1. Configuração da conexão com o MongoDB
	client, cancel, err := connectClient("mongodb://10.101.161.51:27017")
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	defer cancel()

	// 2. Definindo o PID a ser procurado
	var pidToFind int // Declara a variável para armazenar o PID

	// Pede ao usuário para digitar o PID
	fmt.Print("Por favor, digite o PID que você deseja buscar: ")

	// Espera o usuário digitar o valor e pressionar Enter
	_, err = fmt.Scanln(&pidToFind)
	if err != nil {
		log.Fatalf("Erro ao ler o PID: %v", err)
	}

	fmt.Printf("Buscando o PID: %d\n\n", pidToFind)

	// 3. Verificando a primeira coleção
	fmt.Println("🔍 Procurando na coleção 'dbProcessos.gestaoProcessos'...")
	result1 := findByPID(context.Background(), client, "dbProcessos", "gestaoProcessos", "processo.pid", pidToFind)
	if result1.Error != nil {
		log.Fatalf("Erro ao buscar na coleção 'gestaoProcessos': %v", result1.Error)
	}

	if result1.Found {
		fmt.Printf("✅ Documento com PID %d ENCONTRADO na coleção '%s'.\n", pidToFind, result1.CollectionName)
		prettyPrintJSON(result1.Document)
	} else {
		fmt.Printf("❌ Documento com PID %d NÃO encontrado na coleção '%s'.\n", pidToFind, result1.CollectionName)
	}

	fmt.Println()

	// 4. Verificando a segunda coleção
	fmt.Println("🔍 Procurando na coleção 'dbProcessos.collProcessos360'...")
	result2 := findByPID(context.Background(), client, "dbProcessos", "collProcessos360", "procAdministrativo.dadosGerais.pid", pidToFind)
	if result2.Error != nil {
		log.Fatalf("Erro ao buscar na coleção 'collProcessos360': %v", result2.Error)
	}

	if result2.Found {
		fmt.Printf("✅ Documento com PID %d ENCONTRADO na coleção '%s'.\n", pidToFind, result2.CollectionName)
		prettyPrintJSON(result2.Document)
	} else {
		fmt.Printf("❌ Documento com PID %d NÃO encontrado na coleção '%s'.\n", pidToFind, result2.CollectionName)
	}

	// 5. Verificação final
	if result1.Found && result2.Found {
		fmt.Printf("\n✨ O documento com o PID %d foi encontrado em AMBAS as coleções.\n", pidToFind)
	} else if result1.Found {
		fmt.Printf("\n⚠️ O documento com o PID %d foi encontrado apenas na coleção 'gestaoProcessos'.\n", pidToFind)
	} else if result2.Found {
		fmt.Printf("\n⚠️ O documento com o PID %d foi encontrado apenas na coleção 'collProcessos360'.\n", pidToFind)
	} else {
		fmt.Printf("\n🚨 O documento com o PID %d NÃO foi encontrado em nenhuma das coleções.\n", pidToFind)
	}

	fmt.Println("Operação concluída.")
	//fmt.Println("Pressione Enter para sair...")
	//fmt.Scanln()
}
