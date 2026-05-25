package setting

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret          string
	PageSize           int
	PrefixUrl          string
	SnowFlakeMachineId int64

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	UserCenterApiHost string
	SeatApiHost       string
	BpmApiHost        string
	RobotApiHost      string
	HomeHost          string
	ConfigEnv         string
}

type GaiaApi struct {
	BaseUrl      string
	AuthPath     string
	ClientSecret string
	CorpID       string
	TokenTTL     time.Duration

	LeaveSubmitPath    string
	MyApplicationsPath string
	LeaveQuotaPath     string
	LeaveTypesPath     string
	LeaveHoursPath     string
	ExceptionListPath  string
}

var GaiaApiSetting = &GaiaApi{}

// GaiaOpenAPISetting is a backward-compatible alias for previous naming.
var GaiaOpenAPISetting = GaiaApiSetting

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type KafkaSetting struct {
	Host     string
	Username string
	Password string
}

type KafkaTopicSetting struct {
	Topic         string
	ConsumerGroup string
}

var KafkaCommonSetting = &KafkaSetting{}
var BpmTicketOperateKafkaSetting = &KafkaTopicSetting{}
var TicketCanalKafkaSetting = &KafkaTopicSetting{}

type ElasticSetting struct {
	Host     string
	Account  string
	Password string
}

var CommonEsSetting = &ElasticSetting{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	configFile := "conf/app.ini"
	if env := strings.TrimSpace(os.Getenv("ESB_CONFIG_FILE")); env != "" {
		configFile = env
	}
	cfg, err = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, configFile)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse '%s': %v", configFile, err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("elastic", CommonEsSetting)

	mapTo("kafka", KafkaCommonSetting)
	mapTo("kafka_ticket_canal", TicketCanalKafkaSetting)
	mapTo("kafka_bpm_ticket_operate", BpmTicketOperateKafkaSetting)
	mapTo("gaia_api", GaiaApiSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	if GaiaApiSetting.TokenTTL <= 0 {
		GaiaApiSetting.TokenTTL = 3600
	}
	GaiaApiSetting.TokenTTL = GaiaApiSetting.TokenTTL * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

func IsLocalDebug() bool {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip := ipNet.IP.String()
				if ip == "localhost" || ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") {
					return true
				}
			}
		}
	}
	return false
}
