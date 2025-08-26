// Arquivo: database.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SearchResult representa o resultado da busca em uma coleção.
type SearchResult struct {
	CollectionName string
	Document       bson.M
	Found          bool
	Error          error
}

// prettyPrintJSON formata e imprime um documento BSON como JSON.
func prettyPrintJSON(doc bson.M) {
	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		log.Printf("Erro ao converter para JSON: %v", err)
		return
	}
	fmt.Println(string(jsonData))
	fmt.Println("--------------------------------------------------------")
}

// connectClient conecta-se ao MongoDB e retorna o cliente e a função de cancelamento.
func connectClient(uri string) (*mongo.Client, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("erro ao conectar ao MongoDB: %w", err)
	}

	return client, cancel, nil
}

// findByPID busca um documento com um PID específico em uma coleção.
func findByPID(ctx context.Context, client *mongo.Client, dbName, collectionName, pidFieldName string, pid int) SearchResult {
	coll := client.Database(dbName).Collection(collectionName)
	filter := bson.D{{pidFieldName, pid}}

	var doc bson.M
	err := coll.FindOne(ctx, filter).Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return SearchResult{
				CollectionName: collectionName,
				Found:          false,
				Error:          nil,
			}
		}
		return SearchResult{
			CollectionName: collectionName,
			Found:          false,
			Error:          fmt.Errorf("erro na query da coleção '%s': %w", collectionName, err),
		}
	}

	return SearchResult{
		CollectionName: collectionName,
		Document:       doc,
		Found:          true,
		Error:          nil,
	}
}
