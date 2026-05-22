package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AWSS3Buckets struct {
	ProductsBucket string `json:"products_bucket"`
}

type Config struct {
	ServiceName                   string       `json:"service_name"`
	AuthToken                     string       `json:"auth_token"`
	JWTSecret                     string       `json:"jwt_secret"`
	Concurrency                   int          `json:"concurrency"`
	ListenPort                    int          `json:"listen_port"`
	LogLevel                      string       `json:"log_level"`
	AWSS3Region                   string       `json:"aws_s3_region"`
	AWSS3Bucket                   AWSS3Buckets `json:"aws_s3_bucket"`
	BodyUploadLimit               int          `json:"body_upload_limit"`
	ServiceUrl                    string       `json:"service_url"`
	ArtifacthubEsperAppsUrlPrefix string       `json:"artifacthub_products_url_prefix"`
	Db                            DbConfig
	ExternalServices              ExternalServices
}

type ExternalServices struct {
}

type ExternalServiceConfig struct {
	ApiUrl    string `json:"api_url"`
	AuthToken string `json:"auth_token"`
}

var appConfig *Config = nil

func NewConfig() *Config {
	// ensures that only one config object is created across the system
	if appConfig != nil {
		return appConfig
	}

	config := &Config{}
	var ok bool
	var err error

	config.ServiceName = SERVICE_NAME

	listenPortStr, ok := os.LookupEnv("LISTEN_PORT")
	if !ok {
		panic("Missing required variable: LISTEN_PORT")
	}
	config.ListenPort, err = strconv.Atoi(listenPortStr)
	if err != nil {
		panic("Incorrect value type for LISTEN_PORT")
	}

	config.ServiceUrl, ok = os.LookupEnv("SERVICE_URL")
	if !ok {
		config.ServiceUrl = "http://localhost:" + listenPortStr
	}

	// config.AWSS3Region, ok = os.LookupEnv("AWS_S3_REGION")
	// if !ok {
	// 	panic("Missing required variable: AWS_S3_REGION")
	// }

	// config.AWSS3Bucket.ProductsBucket, ok = os.LookupEnv("AWS_S3_PRODUCTS_BUCKET")
	// if !ok {
	// 	panic("Missing required variable: AWS_S3_PRODUCTS_BUCKET")
	// }

	// _, ok = os.LookupEnv("AWS_ACCESS_KEY_ID")
	// if !ok {
	// 	panic("Missing required variable: AWS_ACCESS_KEY_ID")
	// }

	// _, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	// if !ok {
	// 	panic("Missing required variable: AWS_SECRET_ACCESS_KEY")
	// }

	config.LogLevel, ok = os.LookupEnv("LOG_LEVEL")
	if !ok {
		config.LogLevel = "ERROR"
	}

	bodyUploadLimitString, ok := os.LookupEnv("BODY_UPLOAD_LIMIT")
	if !ok {
		config.BodyUploadLimit = (2*1024*1024*1024 - 1) // set body limit as 2GB / int32 max value
	} else {
		config.BodyUploadLimit, err = strconv.Atoi(bodyUploadLimitString)
		if err != nil {
			panic("Incorrect value type for BODY_UPLOAD_LIMIT")
		}
	}

	config.JWTSecret, ok = os.LookupEnv("JWT_SECRET")
	if !ok {
		panic("Missing required variable: JWT_SECRET")
	}

	config.LogLevel = strings.ToUpper(config.LogLevel)
	config.Concurrency = 5
	config.Db = parseDbConfig()
	setExternalService(config)
	appConfig = config
	return appConfig
}

func (c *Config) String() string {
	s := fmt.Sprintf(
		"ServiceName: %s, Concurrency: %d, LOG_LEVEL: %s, Listen Port: %d",
		c.ServiceName, c.Concurrency, c.LogLevel, c.ListenPort,
	)
	return s
}

func setExternalService(config *Config) {
}

func parseDbConfig() DbConfig {
	var err error
	var dbConfig DbConfig
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		panic("Missing required variable: DB_HOST")
	} else {
		dbConfig.Host = dbHost
	}
	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		panic("Missing required variable: DB_PORT")
	} else {
		dbConfig.Port, err = strconv.Atoi(dbPort)
		if err != nil {
			panic("Incorrect value type for DB_PORT")
		}
	}
	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		panic("Missing required variable: DB_USER")
	} else {
		dbConfig.Username = dbUser
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		panic("Missing required variable: DB_PASSWORD")
	} else {
		dbConfig.Password = dbPassword
	}
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		panic("Missing required variable: DB_NAME")
	} else {
		dbConfig.Name = dbName
	}
	return dbConfig
}
