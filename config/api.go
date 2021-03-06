package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mesg-foundation/core/version"
	"github.com/spf13/viper"
)

// All the configuration keys.
const (
	APIServerAddress       = "Api.Server.Address"
	APIServerSocket        = "Api.Server.Socket"
	APIClientTarget        = "Api.Client.Target"
	APIServiceTargetPath   = "Api.Service.TargetPath"
	APIServiceTargetSocket = "Api.Service.TargetSocket"
	APIServiceSocketPath   = "Api.Service.SocketPath"
	LogFormat              = "Log.Format"
	LogLevel               = "Log.Level"
	ServicePathHost        = "Service.Path.Host"
	ServicePathDocker      = "Service.Path.Docker"
	MESGPath               = "MESG.Path"
	CoreImage              = "Core.Image"
)

func setAPIDefault() {
	configPath, _ := getConfigPath()

	viper.SetDefault(MESGPath, configPath)

	viper.SetDefault(APIServerAddress, ":50052")
	viper.SetDefault(APIServerSocket, "/mesg/server.sock")
	os.MkdirAll("/mesg", os.ModePerm)

	viper.SetDefault(APIClientTarget, viper.GetString(APIServerAddress))

	viper.SetDefault(APIServiceSocketPath, filepath.Join(viper.GetString(MESGPath), "server.sock"))
	viper.SetDefault(APIServiceTargetPath, "/mesg/server.sock")
	viper.SetDefault(APIServiceTargetSocket, "unix://"+viper.GetString(APIServiceTargetPath))

	viper.SetDefault(LogFormat, "text")
	viper.SetDefault(LogLevel, "info")

	viper.SetDefault(ServicePathHost, filepath.Join(viper.GetString(MESGPath), "services"))
	viper.SetDefault(ServicePathDocker, filepath.Join("/mesg", "services"))
	os.MkdirAll(viper.GetString(ServicePathDocker), os.ModePerm)

	// Keep only the first part if Version contains space
	coreTag := strings.Split(version.Version, " ")
	viper.SetDefault(CoreImage, "mesg/core:"+coreTag[0])
}
