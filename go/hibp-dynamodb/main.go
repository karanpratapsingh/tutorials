package main

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Schema struct {
	Hash  string `dynamodbav:"hash"`
	Times int    `dynamodbav:"times"`
	Type  string `dynamodbav:"type"`
}

var table string = "Hibp"
var dir string = "chunks"

func main() {
	logFile, writer := getLogFile()
	log.SetOutput(writer)
	defer logFile.Close()

	log.Println("Using table", table, "with directory", dir)

	files := getFiles(dir)

	for num, file := range files {
		filename := file.Name()
		path := "chunks/" + filename

		log.Println("====", num+1, "====")
		log.Println("Starting:", filename)

		file, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		items := []dynamodbTypes.WriteRequest{}

		for scanner.Scan() {
			line := scanner.Text()

			schema := parseLine(line)
			attribute := getAttributes(schema)

			item := dynamodbTypes.WriteRequest{
				PutRequest: &dynamodbTypes.PutRequest{
					Item: attribute,
				},
			}

			items = append(items, item)
		}

		chunks := createChunks(items)
		batches := createBatches(chunks)

		log.Println("Created", len(batches), "batches for", len(chunks), "chunks with", len(items), "items")

		var wg sync.WaitGroup

		for index, batch := range batches {
			failed := 0
			log.Println("Processing batch", index+1)
			batchWriteToDB(&wg, batch, &failed)
			log.Println("Completed with", failed, "failures")
			wg.Wait()
		}

		log.Println("Processed", filename)

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Done")
}

func getLogFile() (*os.File, io.Writer) {
	file, err := os.OpenFile("logs/job.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	return file, mw
}

func getDynamoDBClient() dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), 5000)
	}))

	cfg.Region = "us-west-2"

	if err != nil {
		log.Fatal(err)
	}

	return *dynamodb.NewFromConfig(cfg)
}

func getFiles(dir string) []fs.FileInfo {
	files, dirReadErr := ioutil.ReadDir("chunks")

	if dirReadErr != nil {
		panic(dirReadErr)
	}

	return files
}

func parseLine(line string) Schema {
	split := strings.Split(line, ":")

	Hash := split[0]
	Times, _ := strconv.Atoi(split[1])
	Type := "SHA-1"

	return Schema{Hash, Times, Type}
}

func getAttributes(schema Schema) map[string]dynamodbTypes.AttributeValue {
	attribute, err := attributevalue.MarshalMap(schema)

	if err != nil {
		log.Println("Error processing:", schema)
		log.Fatal(err.Error())
	}

	return attribute
}

func batchWriteToDB(wg *sync.WaitGroup, data [][]dynamodbTypes.WriteRequest, failed *int) {
	for _, chunk := range data {
		wg.Add(1)

		go func(chunk []dynamodbTypes.WriteRequest, failed *int) {
			defer wg.Done()
			client := getDynamoDBClient()

			_, err := client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]dynamodbTypes.WriteRequest{
					table: chunk,
				},
			})

			if err != nil {
				*failed += 1
				log.Println(err.Error())
			}
		}(chunk, failed)
	}
}

func createChunks(arr []dynamodbTypes.WriteRequest) [][]dynamodbTypes.WriteRequest {
	var chunks [][]dynamodbTypes.WriteRequest
	var size int = 25

	for i := 0; i < len(arr); i += size {
		end := i + size

		if end > len(arr) {
			end = len(arr)
		}

		chunks = append(chunks, arr[i:end])
	}

	return chunks
}

func createBatches(arr [][]dynamodbTypes.WriteRequest) [][][]dynamodbTypes.WriteRequest {
	var batches [][][]dynamodbTypes.WriteRequest
	var size int = 10000

	for i := 0; i < len(arr); i += size {
		end := i + size

		if end > len(arr) {
			end = len(arr)
		}

		batches = append(batches, arr[i:end])
	}

	return batches
}
