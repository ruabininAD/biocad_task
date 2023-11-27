package dbImplement

import (
	"Biocad2/src/message"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"path/filepath"
)

type MongoDB struct {
	GUIDCollection *mongo.Collection
	FileCollection *mongo.Collection
	client         *mongo.Client
	ctx            context.Context
}
type FileStatus struct {
	FileName string `json:"fileName"`
	Status   string `json:"status"`
}

func (db *MongoDB) Init(config *viper.Viper) {
	db.ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(config.GetString("host") + config.GetString("port") + "/")
	var err error

	db.client, err = mongo.Connect(db.ctx, clientOptions) // client
	if err != nil {
		log.Fatal(err)
	}

	err = db.client.Ping(db.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.FileCollection = db.client.Database(config.GetString("database")).Collection(config.GetString("collection"))
	db.GUIDCollection = db.client.Database(config.GetString("database")).Collection(config.GetString("guid_collection"))

}

func (db *MongoDB) AddFileName(FileName string, state string) error {

	data := FileStatus{FileName: FileName, Status: state}
	log.Printf("data %s FileStatus %s\n", FileName, data)

	messageBSON, err := bson.Marshal(data)
	if err != nil {
		log.Println("Ошибка маршалинга данных message", err)
		return err
	}

	_, err = db.FileCollection.InsertOne(context.Background(), messageBSON)
	if err != nil {
		log.Println("Ошибка вставки данных message", err)
		return err
	}
	return err
}
func (db *MongoDB) addMessage(message message.Message) error {

	messageBSON, err := bson.Marshal(message)
	if err != nil {
		log.Println("Ошибка маршалинга данных message", err)
		return err
	}

	_, err = db.GUIDCollection.InsertOne(context.Background(), messageBSON)
	if err != nil {
		log.Println("Ошибка вставки данных message", err)
		return err
	}

	return nil
}

// Записывает файл из несколльких message в базу данных
func (db *MongoDB) AddFile(messages []message.Message, filePath string) (err error) {

	for _, message := range messages {
		err = db.addMessage(message)
		if err != nil {
			return err
		}
	}

	_, fileName := filepath.Split(filePath)
	status := fmt.Sprint(len(messages)) + " сообщений"
	if err != nil {
		status = fmt.Sprint(err)
	}
	db.AddFileName(fileName, status)

	return err
}

func (db *MongoDB) GetById(UnitGUID string, pageNumber, pageSize int) (data []message.Message, total int64) {

	// Получаем общее количество документов для данного UnitGUID (для дальнейшей пагинации)
	total, err := db.GUIDCollection.CountDocuments(context.Background(), bson.M{"unitguid": UnitGUID})
	if err != nil {
		log.Println(err)
		return nil, total
	}

	// Вычисляем смещение (skip) и лимит (limit) для пагинации
	offset := int64((pageNumber - 1) * pageSize)

	cur, err := db.GUIDCollection.Find(context.Background(), bson.M{"unitguid": UnitGUID}, options.Find().
		SetSkip(offset).
		SetLimit(int64(pageSize)))
	if err != nil {
		return nil, 0
	}

	defer cur.Close(context.Background())

	if err = cur.All(context.Background(), &data); err != nil {
		log.Fatal(err) //fixme
	}

	return data, total
}

func (db *MongoDB) AllGetById(UnitGUID string) (data []message.Message) {

	cur, err := db.GUIDCollection.Find(context.Background(), bson.M{"unitguid": UnitGUID})
	defer cur.Close(context.Background())
	if err != nil {
		return nil
	}

	if err = cur.All(context.Background(), &data); err != nil {
		log.Fatal(err) //fixme
	}

	return data
}
