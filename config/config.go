package config

import (
	"flag"
	"github.com/spf13/viper"
	"log"
)

func Init(ConfigName string) (mainConfig, dbConfig, PDFConfig *viper.Viper) {
	var err error

	mainConfig = viper.New()
	mainConfig.SetConfigName(ConfigName + "Main")
	mainConfig.SetConfigType("yaml")
	mainConfig.AddConfigPath("config/configFile")
	//mainConfig.AddConfigPath("C:\\Users\\ryabi\\GolandProjects\\biocadGo\\config\\configuration")
	err = mainConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}

	dbConfig = viper.New()
	dbConfig.SetConfigName(ConfigName + "dbConfig")
	dbConfig.SetConfigType("yaml")
	dbConfig.AddConfigPath("config/configFile")
	//dbConfig.AddConfigPath("C:\\Users\\ryabi\\GolandProjects\\biocadGo\\config\\configuration")
	err = dbConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}

	PDFConfig = viper.New()
	PDFConfig.SetConfigName(ConfigName + "PDFConfig") // Указываем имя файла без расширения
	PDFConfig.SetConfigType("yaml")
	PDFConfig.AddConfigPath("config/configFile") // Указываем тип файла
	// Указываем путь к файлу, если он в текущем каталоге
	//PDFConfig.AddConfigPath("C:\\Users\\ryabi\\GolandProjects\\biocadGo\\config\\configuration")
	err = PDFConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}

	return mainConfig, dbConfig, PDFConfig
}

func ParseFlags() map[string]string {
	flag.String("config", "dev", "Имя файла для обработки")
	flag.Parse()
	flags := make(map[string]string)

	flag.Visit(func(f *flag.Flag) {
		flags[f.Name] = f.Value.String()
	})
	return flags
}
