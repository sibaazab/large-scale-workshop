package main 
import (     
	"os"
	"log"     
	"github.com/sibaazab/large-scale-workshop.git/Config"     
	"github.com/sibaazab/large-scale-workshop.git/utils"     
	"gopkg.in/yaml.v2" 
)

//hello
	func main() {
		// read configuration file from command line argument 
		if len(os.Args) != 2 {
			utils.Logger.Fatal("Expecting exactly one configuration file")
			os.Exit(1) }
		configFile := os.Args[1]
		configData, err := os.ReadFile(configFile) 
		if err != nil {
			log.Fatalf("error reading file: %v", err)
			os.Exit(2) }

			
		var config Config.ConfigBase
		err = yaml.Unmarshal(configData, &config) // parses YAML 
		if err != nil {
			log.Fatalf("error unmarshaling data: %v", err)
			os.Exit(3) }
		switch config.Type { case "TestService":
		utils.Logger.Printf("Loading service type: %v\n", config.Type)
		default:
		utils.Logger.Fatalf("Unknown configuration type: %v", config.Type)
		os.Exit(4)
		}}