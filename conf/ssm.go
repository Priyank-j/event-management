package conf

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v3"
)

const ssmBatchSize = 10

var ssmMap = preloadSSMKeys(ENV) // stores pre loaded items from ssm store
var svc = ssm.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion(REGION))

type keyStruct struct {
	Name             string   `yaml:"name"`
	Key              string   `yaml:"key"`
	EnvironmentIn    []string `yaml:"environmentIn"`
	EnvironmentNotIn []string `yaml:"environmentNotIn"`
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func matchEnv(input string, pattern string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(input, pattern[:len(pattern)-1])
	}
	return input == pattern
}

func preloadSSMKeys(env string) map[string]string {
	if (ENV == ENV_LOCAL || ENV == "") && LocalParameterStore != devSSM {
		// ignore ssm for local
		return map[string]string{}
	}
	// read from yaml
	b, err := os.ReadFile("config/keys.yaml")
	if err != nil {
		panic("couldn't read from keys.yaml: " + err.Error())
	}
	// process file contents
	keyObjs := []keyStruct{}
	err = yaml.Unmarshal(b, &keyObjs)
	if err != nil {
		panic("couldn't unmarshal from keys.yaml: " + err.Error())
	}
	keys := []string{} // final keys to load from SSM
	// save in map after iterating
	for _, key := range keyObjs {
		save := true // keep default as true
		if len(key.EnvironmentIn) > 0 {
			// check for environmentIn
			save = false
			for _, envPattern := range key.EnvironmentIn {
				if matchEnv(env, envPattern) {
					save = true
					break
				}
			}
		} else if len(key.EnvironmentNotIn) > 0 {
			// check for environmentNotIn
			save = true
			for _, envPattern := range key.EnvironmentNotIn {
				if matchEnv(env, envPattern) {
					save = false
					break
				}
			}
		}
		if !save {
			// skip
			continue
		}
		// save in keys
		keys = append(keys, key.Key)
	}
	// set key objs to nil
	keyObjs = nil
	// load keys from SSM in batches of 10
	batchCount := int(len(keys) / ssmBatchSize)
	if len(keys)%ssmBatchSize != 0 {
		batchCount++
	}
	ssmKeys := make(map[string]string)
	for i := 0; i < batchCount; i++ {
		keysInBatch := keys[i*ssmBatchSize : min(len(keys), (i+1)*ssmBatchSize)]
		params, err := svc.GetParameters(&ssm.GetParametersInput{Names: aws.StringSlice(keysInBatch), WithDecryption: aws.Bool(true)})
		if err != nil {
			panic(fmt.Errorf("error getting parameter: %s", err.Error()))
		}
		for _, param := range params.Parameters {
			ssmKeys[*param.Name] = *param.Value
		}
	}
	// set keys to nil
	keys = nil
	// return the map
	return ssmKeys
}

const productName = "event_management_app"

// TODO: replace the usage of below function by exposing a direct function from conf and using it
// across all places
func getSSMKey(key string) string {
	if (ENV == ENV_LOCAL || ENV == "") && LocalParameterStore != devSSM {
		// ignore ssm for local
		return ""
	}
	value, found := ssmMap[key]
	if !found {
		panic(fmt.Errorf("parameter %s not found in pre loaded SSM Keys", key))
	}
	return value
}
