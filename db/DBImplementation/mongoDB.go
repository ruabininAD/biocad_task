package DBImplementation

import (
	"biocadGo/src/message"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
)

type MongoDB struct {
	collection *mongo.Collection
	client     *mongo.Client
	ctx        context.Context
}

func (db *MongoDB) Init(host, port string) {
	db.ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(host + port + "/")
	var err error

	db.client, err = mongo.Connect(db.ctx, clientOptions) // client
	if err != nil {
		log.Fatal(err)
	}
	err = db.client.Ping(db.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.ToCollection(viper.GetString("database"), viper.GetString("collection"))
}

func (db *MongoDB) ToCollection(Database, Collection string) {
	db.collection = db.client.Database(Database).Collection(Collection)
}

//// Пример добавления одного документа
//doc := bson.D{{"name", "Aliceeeeeeeeeee"}, {"age", 30}}
//_, err = collection.InsertOne(context.Background(), doc)
//if err != nil {
//	panic(err)
//}

// Пример добавления нескольких документов
//docs := []interface{}{
//	bson.D{{"name", "Bob"}, {"age", 35}},
//	bson.D{{"name", "Charlie"}, {"age", 40}},
//}
//_, err = collection.InsertMany(context.Background(), docs)
//if err != nil {
//	panic(err)
//}
//}

func (db *MongoDB) addMessage(message message.Message) error {

	messageBSON, err := bson.Marshal(message)
	if err != nil {
		fmt.Println("Ошибка маршалинга данных message", err)
		return err
	}

	_, err = db.collection.InsertOne(context.Background(), messageBSON)
	if err != nil {
		fmt.Println("Ошибка вставки данных message", err)
		return err
	}

	return nil
}

// Записывает файл из несколльких message в базу данных
func (db *MongoDB) AddFile(messages []message.Message) error {

	for _, message := range messages {
		err := db.addMessage(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *MongoDB) GetById(UnitGUID string, pageNumber, pageSize int) (data []message.Message, total int64) {

	// Получаем общее количество документов для данного UnitGUID (для дальнейшей пагинации)
	total, err := db.collection.CountDocuments(context.Background(), bson.M{"unitguid": UnitGUID})
	if err != nil {
		fmt.Println(err)
		return nil, total
	}

	// Вычисляем смещение (skip) и лимит (limit) для пагинации
	offset := int64((pageNumber - 1) * pageSize)

	cur, err := db.collection.Find(context.Background(), bson.M{"unitguid": UnitGUID}, options.Find().
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
